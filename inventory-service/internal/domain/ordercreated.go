package domain

type OrderCreated struct {
	EventId string `json:"event_id"` // uuid
	OrderId uint   `json:"order_id"`

	Items []Item `json:"items"`
}

type Item struct {
	ProductID uint `json:"product_id"`
	Qty       uint `json:"qty"`
}
