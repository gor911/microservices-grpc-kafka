package publisher

import (
	"context"
	"time"

	"github.com/gor911/microservices-grpc-kafka/order-service/internal/service"
)

type Publisher struct {
	outboxService *service.Outbox
}

func New(o *service.Outbox) *Publisher {
	return &Publisher{outboxService: o}
}

func (p *Publisher) Run(ctx context.Context) error {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for range t.C {
		if err := p.outboxService.Produce(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (p *Publisher) Close(ctx context.Context) error {
	return nil
}
