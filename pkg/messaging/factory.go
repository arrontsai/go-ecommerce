package messaging

import (
	"context"
	"fmt"
	"strings"

	"github.com/yourusername/ecommerce/pkg/config"
	"go.uber.org/zap"
	"github.com/segmentio/kafka-go"
	"github.com/google/uuid"
)

// MessageBrokerType 表示消息代理類型
type MessageBrokerType string

const (
	// KafkaBroker 表示 Kafka 消息代理
	KafkaBroker MessageBrokerType = "kafka"
)

// MessageBrokerFactory 是消息代理工廠
type MessageBrokerFactory struct {
	Config *config.Config
	Logger *zap.Logger
}

// NewMessageBrokerFactory 創建一個新的消息代理工廠
func NewMessageBrokerFactory(cfg *config.Config, logger *zap.Logger) *MessageBrokerFactory {
	return &MessageBrokerFactory{
		Config: cfg,
		Logger: logger,
	}
}

// KafkaClient 表示 Kafka 客戶端
type KafkaClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

// CreateKafkaClient 創建一個 Kafka 客戶端
func (f *MessageBrokerFactory) CreateKafkaClient() (*KafkaClient, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  f.Config.KafkaBrokers,
		Topic:    "", // 動態指定
		Balancer: &kafka.Hash{},
	})

	return &KafkaClient{writer: writer}, nil
}

// Publish 發送訊息到指定主題
func (k *KafkaClient) Publish(topic string, message []byte) error {
	return k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Key:   []byte(uuid.New().String()),
			Value: message,
		},
	)
}

// CreateMessageBroker 根據類型創建消息代理
func (f *MessageBrokerFactory) CreateMessageBroker(brokerType MessageBrokerType) (interface{}, error) {
	switch brokerType {
	case KafkaBroker:
		return f.CreateKafkaClient()
	default:
		return nil, fmt.Errorf("不支持的消息代理類型: %s", brokerType)
	}
}

// GetPreferredBrokerType 獲取首選的消息代理類型
func (f *MessageBrokerFactory) GetPreferredBrokerType() MessageBrokerType {
	// 預設使用 Kafka
	return KafkaBroker
}
