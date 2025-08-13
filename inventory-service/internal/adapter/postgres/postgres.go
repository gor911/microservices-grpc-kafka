package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/domain"
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

	err = db.AutoMigrate(&domain.Product{}, &domain.Inventory{})

	if err != nil {
		return nil, err
	}

	seedIfEmpty(db)

	return &Postgres{gormDB: db}, nil
}

func seedIfEmpty(gdb *gorm.DB) {
	var cnt int64

	if err := gdb.Model(&domain.Product{}).Count(&cnt).Error; err != nil {
		log.Printf("seed: count error: %v", err)
		return
	}

	if cnt > 0 {
		log.Printf("seed: products already exist (%d), skip seeding", cnt)
		return
	}

	products := []domain.Product{
		{ID: 1, Name: "Keyboard", Price: 1999},
		{ID: 2, Name: "Mouse", Price: 1499},
		{ID: 3, Name: "Case", Price: 3599},
	}

	inv := []domain.Inventory{
		{ProductID: 1, AvailableQty: 200},
		{ProductID: 2, AvailableQty: 150},
		{ProductID: 3, AvailableQty: 100},
	}

	if err := gdb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&products).Error; err != nil {
			return err
		}

		if err := tx.Create(&inv).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Printf("seed: tx error: %v", err)

		return
	}

	log.Printf("seed: inserted %d products and %d inventory rows", len(products), len(inv))
}

func (p *Postgres) Close(ctx context.Context) error {
	db, err := p.gormDB.DB()

	if err != nil {
		return err
	}

	return db.Close()
}

func (p *Postgres) GetProducts(ctx context.Context, ids []uint) ([]*domain.Product, error) {
	var products []*domain.Product

	if err := p.gormDB.Find(&products, ids).Error; err != nil {
		return nil, err
	}

	return products, nil
}
