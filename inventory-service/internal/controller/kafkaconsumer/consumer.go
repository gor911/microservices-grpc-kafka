package kafkaconsumer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string
	Topic   string
}

type Handler func(ctx context.Context, msg kafka.Message) error

type Consumer struct {
	config   Config
	consumer *kafka.Reader
	logger   *slog.Logger
	handler  Handler
}

func NewConsumer(config Config, l *slog.Logger, h Handler) *Consumer {
	return &Consumer{
		config: config,
		consumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  config.Brokers,
			Topic:    config.Topic,
			GroupID:  "order-service-group",
			MaxBytes: 10e6, // 10MB
		}),
		logger:  l,
		handler: h,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		m, err := c.consumer.FetchMessage(ctx)

		if err != nil {
			return fmt.Errorf("consumer error: %w", err)
		}

		if err := c.handler(ctx, m); err != nil {
			return fmt.Errorf("handler error: %w", err)
		}

		if err := c.consumer.CommitMessages(ctx, m); err != nil {
			return fmt.Errorf("failed to commit messages: %w", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
