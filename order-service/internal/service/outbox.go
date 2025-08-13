package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/kafkaproducer"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/postgres"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	Produce(ctx context.Context, messages ...kafka.Message) error
}

type Outbox struct {
	postgres      Postgres
	logger        *slog.Logger
	kafkaProducer KafkaProducer
}

func NewOutbox(p *postgres.Postgres, l *slog.Logger, kp *kafkaproducer.Producer) *Outbox {
	return &Outbox{postgres: p, logger: l, kafkaProducer: kp}
}

func (p *Outbox) Produce(ctx context.Context) error {
	outboxes, err := p.postgres.GetOutboxes(ctx)

	if err != nil {
		return fmt.Errorf("get outboxes: %w", err)
	}

	if len(outboxes) == 0 {
		return nil
	}

	var messages []kafka.Message
	var outboxIds []uint

	for _, outbox := range outboxes {
		key := fmt.Sprintf("%d", outbox.Payload.OrderId)

		payloadBytes, err := json.Marshal(outbox.Payload)

		if err != nil {
			return fmt.Errorf("marshal payload: %w", err)
		}

		messages = append(messages, kafka.Message{Topic: outbox.Topic, Key: []byte(key), Value: payloadBytes})
		outboxIds = append(outboxIds, outbox.ID)
	}

	if err := p.kafkaProducer.Produce(ctx, messages...); err != nil {
		return fmt.Errorf("produce outboxes: %w", err)
	}

	if err := p.postgres.MarkOutboxesAsSent(ctx, outboxIds); err != nil {
		return fmt.Errorf("mark outboxes: %w", err)
	}

	return nil
}
