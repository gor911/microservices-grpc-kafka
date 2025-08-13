package domain

import (
	"time"

	"github.com/gor911/microservices-grpc-kafka/order-service/internal/dto"
	invpb "github.com/gor911/microservices-grpc-kafka/order-service/pkg/grpc/inventorypb"
)

type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	Total     uint64
	Status    string // PENDING, CONFIRMED, CANCELLED
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	ProductId uint
	Qty       uint
	Price     uint64
	Name      string
}

func NewOrder(input dto.CreateOrderInput, products []*invpb.Product) *Order {
	o := &Order{}
	o.UserID = input.UserID
	o.Status = "PENDING"

	qtyPerIds := make(map[uint]uint, len(input.Items))

	for _, item := range input.Items {
		qtyPerIds[item.ProductID] = item.Qty
	}

	for _, product := range products {
		o.Items = append(o.Items, OrderItem{
			ProductId: uint(product.Id),
			Name:      product.Name,
			Price:     product.Price,
			Qty:       qtyPerIds[uint(product.Id)],
		})
	}

	return o
}
