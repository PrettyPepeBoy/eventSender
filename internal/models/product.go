package models

type Product struct {
	Name     string `json:"product_name" validate:"required"`
	Category string `json:"product_category" validate:"required"`
}

type ProductWithId struct {
	ProductId string `json:"product_id" validate:"required"`
}

type ProductName struct {
	Name string `json:"product_name" validate:"required"`
}
