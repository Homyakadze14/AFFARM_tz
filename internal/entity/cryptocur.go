package entity

import "time"

type Cryptocurrencie struct {
	ID     int
	Symbol string
	Name   string
}

type Tracking struct {
	ID                int
	CryptocurrencieID int
	IsActive          bool
}

type PriceHistory struct {
	ID                int
	CryptocurrencieID int
	Price             float64
	Timestamp         time.Time
}
