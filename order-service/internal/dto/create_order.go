package dto

type CreateOrderInput struct {
	UserID uint `json:"user_id"`

	Items []CreateOrderItemInput `json:"items"`
}

type CreateOrderItemInput struct {
	ProductID uint `json:"product_id"`
	Qty       uint `json:"qty"`
}

type CreateOrderOutput struct {
	OrderID uint `json:"order_id"`
}
