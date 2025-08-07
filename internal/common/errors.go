package common

import "errors"

var (
	ErrCryptocurrencyAlreadyExists = errors.New("cryptocurrency already exists")
	ErrCryptocurrencyNotFound      = errors.New("cryptocurrency not found")
	ErrTrackingAlreadyExists       = errors.New("tracking already exists")
	ErrTrackingNotFound            = errors.New("tracking not found")
	ErrHistoryNotFound             = errors.New("history not found")
)
