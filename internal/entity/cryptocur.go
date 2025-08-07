package entity

import "time"

type Cryptocurrency struct {
	ID     int
	Symbol string
}

type Tracking struct {
	ID               int
	CryptocurrencyID int
	IsActive         bool
}

type PriceHistory struct {
	ID               int
	CryptocurrencyID int
	Price            float64
	Timestamp        time.Time
}
