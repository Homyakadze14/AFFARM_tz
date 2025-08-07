package dto

type AddCryptocurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required"`
}

type RemoveCryptocurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required"`
}

type PriceRequest struct {
	Symbol    string `json:"symbol" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

type PriceResponse struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}
