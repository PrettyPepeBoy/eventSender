package models

type User struct {
	Id   int64  `json:"id,omitempty"`
	Mail string `json:"mail" validate:"required,email"`
}
