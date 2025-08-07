package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Homyakadze14/AFFARM_tz/internal/common"
	"github.com/Homyakadze14/AFFARM_tz/internal/entity"
)

type CryptocurrencyStorage interface {
	Create(ctx context.Context, cryptocur *entity.Cryptocurrency) (*entity.Cryptocurrency, error)
	GetBySymbol(ctx context.Context, symbol string) (*entity.Cryptocurrency, error)
	GetAll(ctx context.Context) ([]entity.Cryptocurrency, error)
	CreateOrGet(ctx context.Context, c *entity.Cryptocurrency) (*entity.Cryptocurrency, error)
}

type TrackingStorage interface {
	Create(ctx context.Context, trc *entity.Tracking) (*entity.Tracking, error)
	GetByCryptocurrencyID(ctx context.Context, crid int) (*entity.Tracking, error)
	GetActive(ctx context.Context) ([]entity.Tracking, error)
	Update(ctx context.Context, trc *entity.Tracking) (*entity.Tracking, error)
}

type HistoryStorage interface {
	Create(ctx context.Context, history *entity.PriceHistory) (*entity.PriceHistory, error)
	GetNearestPrice(ctx context.Context, cryptocurrencyID int, timestamp time.Time) (*entity.PriceHistory, error)
}

type CryptoClient interface {
	GetPrice(symbol string, currency string) (float64, error)
	SymbolExists(symbol string) (bool, error)
}

type CryptocurrencyService struct {
	log         *slog.Logger
	cst         CryptocurrencyStorage
	tst         TrackingStorage
	hst         HistoryStorage
	cryptoCient CryptoClient
}

func NewCryptocurrencyService(
	log *slog.Logger,
	cst CryptocurrencyStorage,
	tst TrackingStorage,
	hst HistoryStorage,
	cryptoCient CryptoClient,
) *CryptocurrencyService {
	return &CryptocurrencyService{
		cst:         cst,
		tst:         tst,
		hst:         hst,
		cryptoCient: cryptoCient,
	}
}

func (s *CryptocurrencyService) Add(ctx context.Context, cr *entity.Cryptocurrency) error {
	const op = "CryptocurrencyService.Add"
	log := s.log.With(slog.String("op", op),
		slog.String("symbol", cr.Symbol),
		slog.String("name", cr.Name))

	log.Debug("trying to add cryptocurrency")
	exists, err := s.cryptoCient.SymbolExists(cr.Symbol)
	if err != nil {
		log.Error(fmt.Sprintf("fail to check existance! Error: %s", err))
		return err
	}

	if !exists {
		log.Error("symbol doesn't exists!")
		return common.ErrSymbolNotFound
	}

	cr, err = s.cst.CreateOrGet(ctx, cr)
	if err != nil {
		log.Error(fmt.Sprintf("fail to create or get cryptocurrency! Error: %s", err))
		return err
	}

	trc, err := s.tst.GetByCryptocurrencyID(ctx, cr.ID)
	if err != nil {
		if !errors.Is(err, common.ErrTrackingNotFound) {
			log.Error(fmt.Sprintf("fail to get tracking! Error: %s", err))
			return err
		}

		tr := &entity.Tracking{
			CryptocurrencyID: cr.ID,
			IsActive:         true,
		}
		trc, err = s.tst.Create(ctx, tr)
		if err != nil {
			log.Error(fmt.Sprintf("fail to create tracking! Error: %s", err))
			return err
		}

		// TODO: Запуск сбора цены
	}

	if !trc.IsActive {
		trc.IsActive = true
		trc, err = s.tst.Update(ctx, trc)
		if err != nil {
			log.Error(fmt.Sprintf("fail to update tracking! Error: %s", err))
			return err
		}

		// TODO: Запуск сбора цены
	}
	log.Debug("successfully added cryptocurrency")

	return nil
}

func (s *CryptocurrencyService) Remove(ctx context.Context, cr *entity.Cryptocurrency) error {
	const op = "CryptocurrencyService.Remove"
	log := s.log.With(slog.String("op", op),
		slog.String("symbol", cr.Symbol))

	log.Debug("trying to remove cryptocurrency")
	cr, err := s.cst.GetBySymbol(ctx, cr.Symbol)
	if err != nil {
		log.Error(fmt.Sprintf("fail to get cryptocurrency by symbol! Error: %s", err))
		return err
	}

	tr, err := s.tst.GetByCryptocurrencyID(ctx, cr.ID)
	if err != nil {
		log.Error(fmt.Sprintf("fail to get tracking by cryptocurrency! Error: %s", err))
		return err
	}

	tr.IsActive = false
	_, err = s.tst.Update(ctx, tr)
	if err != nil {
		log.Error(fmt.Sprintf("fail to update tracking! Error: %s", err))
		return err
	}

	// TODO: Удалить задачу сбора цен
	log.Debug("successfully removed cryptocurrency")

	return nil
}

func (s *CryptocurrencyService) Price(ctx context.Context, symbol string, timestamp time.Time) (*entity.PriceHistory, error) {
	const op = "CryptocurrencyService.Price"
	log := s.log.With(slog.String("op", op),
		slog.String("symbol", symbol),
		slog.Time("timestamp", timestamp))

	log.Debug("trying to get price of cryptocurrency")
	cr, err := s.cst.GetBySymbol(ctx, symbol)
	if err != nil {
		log.Error(fmt.Sprintf("fail to get cryptocurrency by symbol! Error: %s", err))
		return nil, err
	}

	hist, err := s.hst.GetNearestPrice(ctx, cr.ID, timestamp)
	if err != nil {
		log.Error(fmt.Sprintf("fail to get nearest price! Error: %s", err))
		return nil, err
	}
	log.Debug("successfully got price of cryptocurrency")

	return hist, nil
}
