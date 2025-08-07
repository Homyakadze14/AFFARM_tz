package common

import "errors"

var (
	ErrCryptocurrencyAlreadyExists = errors.New("cryptocurrency already exists")
	ErrCryptocurrencyNotFound      = errors.New("cryptocurrency not found")
)
