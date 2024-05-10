package storage

import "errors"

var (
	ErrUserAlreadyExist = errors.New("user is already exist")
	ErrUserNotExist     = errors.New("user is not exist")
)
