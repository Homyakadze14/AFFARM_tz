package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Homyakadze14/AFFARM_tz/internal/common"
	"github.com/Homyakadze14/AFFARM_tz/internal/entity"
	"github.com/Homyakadze14/AFFARM_tz/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type HistoryRepo struct {
	*postgres.Postgres
}

func NewHistoryRepository(pg *postgres.Postgres) *HistoryRepo {
	return &HistoryRepo{pg}
}

func (r *HistoryRepo) Create(ctx context.Context, history *entity.PriceHistory) (*entity.PriceHistory, error) {
	const op = "HistoryRepo.Create"

	err := r.Pool.QueryRow(ctx,
		`INSERT INTO price_history (cryptocurrency_id, price, timestamp)
		VALUES ($1, $2, $3)
		RETURNING id;`, history.CryptocurrencyID, history.Price, history.Timestamp).Scan(&history.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return history, nil
}

func (r *HistoryRepo) GetNearestPrice(ctx context.Context, cryptocurrencyID int, timestamp time.Time) (*entity.PriceHistory, error) {
	const op = "HistoryRepo.GetNearestPrice"

	// 1. Пробуем найти точное совпадение
	exact, err := r.getExactPrice(ctx, cryptocurrencyID, timestamp)
	if err == nil {
		return exact, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// 2. Ищем ближайшую запись до указанного времени
	before, err := r.getNearestBefore(ctx, cryptocurrencyID, timestamp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// 3. Ищем ближайшую запись после указанного времени
	after, err := r.getNearestAfter(ctx, cryptocurrencyID, timestamp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// 4. Выбираем ближайший результат
	switch {
	case before == nil && after == nil:
		return nil, common.ErrHistoryNotFound
	case before == nil:
		return after, nil
	case after == nil:
		return before, nil
	default:
		beforeDiff := timestamp.Sub(before.Timestamp)
		afterDiff := after.Timestamp.Sub(timestamp)

		if beforeDiff <= afterDiff {
			return before, nil
		}
		return after, nil
	}
}

func (r *HistoryRepo) getExactPrice(ctx context.Context, cryptocurrencyID int, timestamp time.Time) (*entity.PriceHistory, error) {
	const op = "HistoryRepo.getExactPrice"

	const query = `SELECT id, cryptocurrency_id, price, timestamp 
                   FROM price_history 
                   WHERE cryptocurrency_id = $1 AND timestamp = $2`

	var history entity.PriceHistory
	err := r.Pool.QueryRow(ctx, query, cryptocurrencyID, timestamp).Scan(
		&history.ID,
		&history.CryptocurrencyID,
		&history.Price,
		&history.Timestamp,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &history, nil
}

func (r *HistoryRepo) getNearestBefore(ctx context.Context, cryptocurrencyID int, timestamp time.Time) (*entity.PriceHistory, error) {
	const op = "HistoryRepo.getNearestBefore"
	const query = `SELECT id, cryptocurrency_id, price, timestamp 
                   FROM price_history 
                   WHERE cryptocurrency_id = $1 AND timestamp <= $2
                   ORDER BY timestamp DESC 
                   LIMIT 1`

	var history entity.PriceHistory
	err := r.Pool.QueryRow(ctx, query, cryptocurrencyID, timestamp).Scan(
		&history.ID,
		&history.CryptocurrencyID,
		&history.Price,
		&history.Timestamp,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &history, nil
}

func (r *HistoryRepo) getNearestAfter(ctx context.Context, cryptocurrencyID int, timestamp time.Time) (*entity.PriceHistory, error) {
	const op = "HistoryRepo.getNearestAfter"
	const query = `SELECT id, cryptocurrency_id, price, timestamp 
                   FROM price_history 
                   WHERE cryptocurrency_id = $1 AND timestamp > $2
                   ORDER BY timestamp ASC 
                   LIMIT 1`

	var history entity.PriceHistory
	err := r.Pool.QueryRow(ctx, query, cryptocurrencyID, timestamp).Scan(
		&history.ID,
		&history.CryptocurrencyID,
		&history.Price,
		&history.Timestamp,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &history, nil
}
