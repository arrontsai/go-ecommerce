package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// KafkaClient 是 Kafka 客戶端
type KafkaClient struct {
	Writer     *kafka.Writer
	Readers    map[string]*kafka.Reader
	Logger     *zap.Logger
	BrokerURLs []string
}

// KafkaConfig 是 Kafka 配置
type KafkaConfig struct {
	BrokerURLs []string
	ClientID   string
	GroupID    string
}

// MessageHandler 是消息處理函數類型
type MessageHandler func([]byte) error

// NewKafkaClient 創建一個新的 Kafka 客戶端
func NewKafkaClient(config KafkaConfig, logger *zap.Logger) (*KafkaClient, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger 不能為空")
	}

	if len(config.BrokerURLs) == 0 {
		return nil, fmt.Errorf("broker URLs 不能為空")
	}

	client := &KafkaClient{
		BrokerURLs: config.BrokerURLs,
		Logger:     logger,
		Readers:    make(map[string]*kafka.Reader),
	}

	// 創建一個默認的寫入器
	client.Writer = &kafka.Writer{
		Addr:         kafka.TCP(config.BrokerURLs...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	return client, nil
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

// ConsumeMessages 消費指定主題的消息
func (c *KafkaClient) ConsumeMessages(ctx context.Context, topic, groupID string, handler MessageHandler) error {
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
