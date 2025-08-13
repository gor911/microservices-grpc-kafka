package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/dto"
)

type Outbox struct {
	ID        uint    `gorm:"primaryKey"`
	Topic     string  `gorm:"type:string"`
	Payload   Payload `gorm:"serializer:json"`
	CreatedAt time.Time
	SentAt    sql.NullTime
}

type Payload struct {
	EventId string `json:"event_id"` // uuid
	OrderId uint   `json:"order_id"`

	Items []Item `json:"items"`
}

type Item struct {
	ProductID uint `json:"product_id"`
	Qty       uint `json:"qty"`
}

func NewOutbox(topic string, input dto.CreateOrderInput) *Outbox {
	o := &Outbox{}

	o.Topic = topic
	o.Payload.EventId = uuid.NewString()

	for _, item := range input.Items {
		o.Payload.Items = append(o.Payload.Items, Item{ProductID: item.ProductID, Qty: item.Qty})
	}

	return o
}
