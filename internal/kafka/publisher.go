package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
	topic  string
	logger *log.Logger
}

type FileMetadata struct {
	Bucket     string            `json:"bucket"`
	ObjectName string            `json:"object_name"`
	Tags       map[string]string `json:"tags"`
}

func NewPublisher(brokers []string, topic string, logger *log.Logger) (*Publisher, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	// Можно добавить здесь логику создания топика, если он не существует
	// (в зависимости от прав в Kafka)

	return &Publisher{
		writer: writer,
		topic:  topic,
		logger: logger,
	}, nil
}

func (p *Publisher) Send(ctx context.Context, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx,
		kafka.Message{Value: data},
	)
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
