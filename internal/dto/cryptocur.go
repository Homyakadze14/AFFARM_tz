package dto

import "time"

type AddCryptocurrencyRequest struct {
	Symbol string `json:"symbol"`
}

type RemoveCryptocurrencyRequest struct {
	Symbol string `json:"symbol"`
}

type PriceRequest struct {
	Symbol    string    `json:"symbol"`
	Timestamp time.Time `json:"timestamp"`
}

type PriceResponse struct {
	Price float64 `json:"price"`
}
