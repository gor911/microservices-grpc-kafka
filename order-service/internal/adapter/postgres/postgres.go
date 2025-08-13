package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/gor911/microservices-grpc-kafka/order-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type Postgres struct {
	gormDB *gorm.DB
}

func New(ctx context.Context, config Config) (*Postgres, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=UTC", config.Username, config.Password, config.Database, config.Host, config.Port),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}, &domain.Outbox{})

	if err != nil {
		return nil, err
	}

	return &Postgres{gormDB: db}, nil
}

func (p *Postgres) Close(ctx context.Context) error {
	db, err := p.gormDB.DB()

	if err != nil {
		return err
	}

	return db.Close()
}

func (p *Postgres) CreateOrder(ctx context.Context, order *domain.Order, outbox *domain.Outbox) error {
	err := p.gormDB.Transaction(func(tx *gorm.DB) error {
		if err := p.gormDB.Create(order).Error; err != nil {
			return err
		}

		outbox.Payload.OrderId = order.ID

		if err := p.gormDB.Create(outbox).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (p *Postgres) GetOutboxes(ctx context.Context) ([]*domain.Outbox, error) {
	var outboxes []*domain.Outbox

	// limit
	if err := p.gormDB.
		Where("sent_at IS NULL").
		Order("id asc").
		Limit(10).
		Find(&outboxes).Error; err != nil {
		return nil, fmt.Errorf("failed to get outboxes: %w", err)
	}

	return outboxes, nil
}

func (p *Postgres) MarkOutboxesAsSent(ctx context.Context, ids []uint) error {
	if err := p.gormDB.Model(&domain.Outbox{}).Where("id in ?", ids).Updates(map[string]interface{}{"sent_at": time.Now()}).Error; err != nil {
		return fmt.Errorf("mark outboxes as sent: %w", err)
	}

	return nil
}
