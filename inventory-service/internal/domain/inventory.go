package domain

type Inventory struct {
	ProductID    uint `gorm:"primaryKey"`
	AvailableQty uint
}

//type ProcessedEvent struct {
//	ID          string `gorm:"primaryKey"`
//	Topic       string
//	ProcessedAt time.Time
//}
