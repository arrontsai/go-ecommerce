package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// KafkaProducer 生產者介面
type KafkaProducer interface {
	Produce(topic string, message interface{}) error
}

// KafkaClient 是 Kafka 客戶端
type KafkaClient struct {
	Writer     *kafka.Writer
	Readers    map[string]*kafka.Reader
	Logger     *zap.Logger
	BrokerURLs []string
}

// NewKafkaClient 創建新的 Kafka 客戶端 (直接初始化版本)
func NewKafkaClient(brokers []string, logger *zap.Logger) (*KafkaClient, error) {
	if logger == nil {
		return nil, errors.New("logger 參數不可為空")
	}
	if len(brokers) == 0 {
		return nil, errors.New("至少需要指定一個 Kafka broker")
	}

	return &KafkaClient{
		BrokerURLs: brokers,
		Logger:     logger,
		Readers:    make(map[string]*kafka.Reader),
		Writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: "", // 動態指定
		},
	}, nil
}

// PublishMessage 發布消息到指定的主題
func (c *KafkaClient) PublishMessage(ctx context.Context, topic string, key string, message []byte) error {
	c.Logger.Info("發布消息到 Kafka",
		zap.String("topic", topic),
		zap.String("key", key),
		zap.Int("message_size", len(message)),
	)

	err := c.Writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: message,
		Time:  time.Now(),
	})

	if err != nil {
		c.Logger.Error("發布消息失敗",
			zap.String("topic", topic),
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("發布消息失敗: %w", err)
	}

	return nil
}

// PublishJSON 將對象序列化為 JSON 並發布到指定的主題
func (c *KafkaClient) PublishJSON(ctx context.Context, topic string, key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化 JSON 失敗: %w", err)
	}

	return c.PublishMessage(ctx, topic, key, jsonData)
}

// Publish 發布消息到指定的主題
func (c *KafkaClient) Publish(topic string, message []byte) error {
	c.Logger.Info("Publishing message", zap.String("topic", topic))
	return c.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Value: message,
		},
	)
}

// ConsumeMessages 消費指定主題的消息
func (c *KafkaClient) ConsumeMessages(ctx context.Context, topic, groupID string, handler func([]byte) error) error {
	// 檢查是否已經有這個主題的讀取器
	readerKey := fmt.Sprintf("%s-%s", topic, groupID)
	reader, exists := c.Readers[readerKey]

	if !exists {
		// 創建新的讀取器
		reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:     c.BrokerURLs,
			Topic:       topic,
			GroupID:     groupID,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			StartOffset: kafka.FirstOffset,
			Logger:      kafka.LoggerFunc(log.Printf),
		})
		c.Readers[readerKey] = reader
	}

	// 啟動一個 goroutine 來消費消息
	go func() {
		defer reader.Close()

		for {
			select {
			case <-ctx.Done():
				c.Logger.Info("停止消費消息", zap.String("topic", topic), zap.String("group", groupID))
				return
			default:
				message, err := reader.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != context.Canceled {
						c.Logger.Error("讀取消息失敗",
							zap.String("topic", topic),
							zap.String("group", groupID),
							zap.Error(err),
						)
					}
					continue
				}

				c.Logger.Info("收到消息",
					zap.String("topic", topic),
					zap.String("group", groupID),
					zap.String("key", string(message.Key)),
					zap.Int("message_size", len(message.Value)),
				)

				// 處理消息
				if err := handler(message.Value); err != nil {
					c.Logger.Error("處理消息失敗",
						zap.String("topic", topic),
						zap.String("group", groupID),
						zap.String("key", string(message.Key)),
						zap.Error(err),
					)
				}
			}
		}
	}()

	return nil
}

// Produce 發布結構化消息到指定的主題 (用於實作 KafkaProducer 介面)
func (c *KafkaClient) Produce(topic string, message interface{}) error {
	// 序列化消息
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失敗: %w", err)
	}

	// 使用空 key 發布消息
	return c.PublishMessage(context.Background(), topic, "", messageBytes)
}

// Close 關閉 Kafka 客戶端
func (c *KafkaClient) Close() error {
	// 關閉寫入器
	if c.Writer != nil {
		if err := c.Writer.Close(); err != nil {
			c.Logger.Error("關閉 Kafka 寫入器失敗", zap.Error(err))
		}
	}

	// 關閉所有讀取器
	for key, reader := range c.Readers {
		if err := reader.Close(); err != nil {
			c.Logger.Error("關閉 Kafka 讀取器失敗", zap.String("reader", key), zap.Error(err))
		}
	}

	return nil
}
