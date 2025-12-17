package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer() (*Producer, error) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           1,
		MaxAttempts:            10,
		AllowAutoTopicCreation: true,
		ReadTimeout:            10 * time.Second,
		WriteTimeout:           10 * time.Second,

		Async: true,
	}

	return &Producer{writer: writer}, nil
}

// Fire & Forget => no wait, runs in goroutine
func (p *Producer) Produce(topic string, msg string) error {
	go func() {
		p.writer.Topic = topic

		_ = p.writer.WriteMessages(
			context.Background(),
			kafka.Message{
				Key:   []byte("key"),
				Value: []byte(msg),
			},
		)
	}()

	return nil // respond immediately
}
