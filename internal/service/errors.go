package service

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
)
