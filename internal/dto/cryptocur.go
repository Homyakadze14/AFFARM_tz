package dto

type AddCryptocurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required" example:"BTC"`
}

type RemoveCryptocurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required" example:"BTC"`
}

type PriceRequest struct {
	Symbol    string `json:"symbol" binding:"required" example:"BTC"`
	Timestamp int64  `json:"timestamp" binding:"required" example:"1754578944"`
}

type PriceResponse struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}
