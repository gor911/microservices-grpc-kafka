package kafkaproducer

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string
}

type Producer struct {
	config   Config
	producer *kafka.Writer
	logger   *slog.Logger
}

func NewProducer(l *slog.Logger, config Config) *Producer {
	return &Producer{
		config: config,
		producer: &kafka.Writer{
			Addr: kafka.TCP(config.Brokers...),
			// NOTE: When Topic is not defined here, each Message must define it instead.
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,
		},
		logger: l,
	}

}

func (p *Producer) Produce(ctx context.Context, messages ...kafka.Message) error {
	err := p.producer.WriteMessages(context.Background(),
		// NOTE: Each Message has Topic defined, otherwise an error is returned.
		messages...,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
