package models

type User struct {
	Id       int64  `json:"user_id,omitempty"`
	Mail     string `json:"user_mail" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserWithId struct {
	UserId int64 `json:"user_id" validate:"required"`
}

type UserWithMail struct {
	Mail string `json:"user_mail" validate:"required,email"`
}

type UserIdAndProductName struct {
	UserWithId
	ProductName
}
