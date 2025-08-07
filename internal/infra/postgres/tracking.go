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

const trackDefaultSliceCap = 50

type TrackingRepo struct {
	*postgres.Postgres
}

func NewTrackingRepository(pg *postgres.Postgres) *TrackingRepo {
	return &TrackingRepo{pg}
}

func (r *TrackingRepo) Create(ctx context.Context, trc *entity.Tracking) (*entity.Tracking, error) {
	const op = "TrackingRepo.Create"

	err := r.Pool.QueryRow(ctx,
		`INSERT INTO trackings (cryptocurrency_id, is_active)
		VALUES ($1, $2)
		RETURNING id;`, trc.CryptocurrencyID, trc.IsActive).Scan(&trc.ID)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, fmt.Errorf("%s: %w", op, common.ErrTrackingAlreadyExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return trc, nil
}

func (r *TrackingRepo) get(ctx context.Context, op string, condition string, args ...interface{}) (*entity.Tracking, error) {
	row := r.Pool.QueryRow(ctx,
		fmt.Sprintf("SELECT id, cryptocurrency_id, is_active FROM trackings WHERE %s", condition),
		args...)

	var t entity.Tracking
	err := row.Scan(&t.ID, &t.CryptocurrencyID, &t.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, common.ErrTrackingNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &t, nil
}

func (r *TrackingRepo) GetByCryptocurrencyID(ctx context.Context, crid int) (*entity.Tracking, error) {
	const op = "TrackingRepo.GetByCryptocurrencyID"
	condition := "cryptocurrency_id=$1"

	return r.get(ctx, op, condition, crid)
}

func (r *TrackingRepo) GetActive(ctx context.Context) ([]entity.Tracking, error) {
	const op = "TrackingRepo.GetActive"

	rows, err := r.Pool.Query(ctx,
		"SELECT id, cryptocurrency_id, is_active FROM trackings WHERE is_active=true")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	trcks := make([]entity.Tracking, 0, trackDefaultSliceCap)
	for rows.Next() {
		var trck entity.Tracking

		err := rows.Scan(
			&trck.ID, &trck.CryptocurrencyID, &trck.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		trcks = append(trcks, trck)
	}

	return trcks, nil
}

func (r *TrackingRepo) Update(ctx context.Context, trc *entity.Tracking) (*entity.Tracking, error) {
	const op = "TrackingRepo.Update"

	_, err := r.Pool.Exec(ctx,
		`UPDATE trackings SET is_active=$1 WHERE id=$2`,
		trc.IsActive, trc.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return trc, nil
}
