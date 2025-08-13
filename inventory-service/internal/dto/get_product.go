package dto

type GetProductsInput struct {
	ProductIds []uint `json:"product_ids"`
}

type GetProductsOutput struct {
	Products []Product `json:"products"`
}

type Product struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price uint64 `json:"price"`
}
