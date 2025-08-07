package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Homyakadze14/AFFARM_tz/internal/common"
	"github.com/Homyakadze14/AFFARM_tz/internal/entity"
	"github.com/Homyakadze14/AFFARM_tz/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

const cryptDefaultSliceCap = 50

type CryptocurRepo struct {
	*postgres.Postgres
}

func NewCryptocurrencyRepository(pg *postgres.Postgres) *CryptocurRepo {
	return &CryptocurRepo{pg}
}

func (r *CryptocurRepo) Create(ctx context.Context, cryptocur *entity.Cryptocurrency) (*entity.Cryptocurrency, error) {
	const op = "CryptocurRepo.Create"

	err := r.Pool.QueryRow(ctx,
		`INSERT INTO cryptocurrencies (symbol)
		VALUES ($1)
		RETURNING id;`, cryptocur.Symbol).Scan(&cryptocur.ID)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, fmt.Errorf("%s: %w", op, common.ErrCryptocurrencyAlreadyExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cryptocur, nil
}

func (r *CryptocurRepo) get(ctx context.Context, op string, condition string, args ...interface{}) (*entity.Cryptocurrency, error) {
	row := r.Pool.QueryRow(ctx,
		fmt.Sprintf("SELECT id, symbol FROM cryptocurrencies WHERE %s", condition),
		args...)

	var c entity.Cryptocurrency
	err := row.Scan(c.ID, c.Symbol)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, common.ErrCryptocurrencyNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &c, nil
}

func (r *CryptocurRepo) GetBySymbol(ctx context.Context, symbol string) (*entity.Cryptocurrency, error) {
	const op = "CryptocurRepo.GetBySymbol"
	condition := "symbol=$1"

	return r.get(ctx, op, condition, symbol)
}

func (r *CryptocurRepo) GetAll(ctx context.Context) ([]entity.Cryptocurrency, error) {
	const op = "CryptocurRepo.GetAll"

	rows, err := r.Pool.Query(ctx,
		"SELECT id, symbol, name FROM cryptocurrencies")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	crypts := make([]entity.Cryptocurrency, 0, cryptDefaultSliceCap)
	for rows.Next() {
		var crypt entity.Cryptocurrency

		err := rows.Scan(
			&crypt.ID, &crypt.Symbol,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		crypts = append(crypts, crypt)
	}

	return crypts, nil
}

func (r *CryptocurRepo) CreateOrGet(ctx context.Context, c *entity.Cryptocurrency) (*entity.Cryptocurrency, error) {
	const op = "CryptocurRepo.CreateOrGet"

	c, err := r.Create(ctx, c)
	if err != nil {
		if !errors.Is(err, common.ErrCryptocurrencyAlreadyExists) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		c, err = r.GetBySymbol(ctx, c.Symbol)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return c, nil
}
