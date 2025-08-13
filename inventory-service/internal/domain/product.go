package domain

type Product struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price uint64
}
