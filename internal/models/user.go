package models

type User struct {
	Id   int64  `json:"id,omitempty"`
	Mail string `json:"mail" validate:"required,email"`
}

type UserWithId struct {
	UserId int64 `json:"user_id" validate:"required"`
}

type ProductWithId struct {
	ProductId int64 `json:"product_id" validate:"required"`
}
type UserAndProductId struct {
	UserId    int64 `json:"user_id"`
	ProductId int64 `json:"product_id"`
}
