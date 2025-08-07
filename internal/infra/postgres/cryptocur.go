package postgres

import (
	"github.com/Homyakadze14/AFFARM_tz/pkg/postgres"
)

type Repo struct {
	*postgres.Postgres
}

func NewCryptocurrencyRepository(pg *postgres.Postgres) *Repo {
	return &Repo{pg}
}
