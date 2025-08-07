package background

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Homyakadze14/AFFARM_tz/internal/entity"
)

const taskChanSize = 50

type HistoryStorage interface {
	Create(ctx context.Context, history *entity.PriceHistory) (*entity.PriceHistory, error)
}

type CurrencyStorage interface {
	GetActive(ctx context.Context) ([]entity.Cryptocurrency, error)
}

type CryptoClient interface {
	GetPrice(symbol string, currency string) (float64, error)
}

type Parser struct {
	coins          sync.Map
	log            *slog.Logger
	updateInterval time.Duration
	maxWorkers     int
	hst            HistoryStorage
	cst            CurrencyStorage
	cryptoClient   CryptoClient
	done           chan struct{}
	wg             sync.WaitGroup
}

func NewParser(
	log *slog.Logger,
	updateInterval time.Duration,
	maxWorkers int,
	hst HistoryStorage,
	cst CurrencyStorage,
	cryptoClient CryptoClient,
) *Parser {
	return &Parser{
		log:            log,
		updateInterval: updateInterval,
		maxWorkers:     maxWorkers,
		hst:            hst,
		cst:            cst,
		cryptoClient:   cryptoClient,
	}
}

func (p *Parser) Start() {
	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	defer done()
	crs, err := p.cst.GetActive(ctx)
	if err != nil {
		panic(err)
	}

	for _, cr := range crs {
		p.coins.Store(cr.Symbol, cr)
	}

	taskChan := make(chan entity.Cryptocurrency, taskChanSize)
	p.done = make(chan struct{})

	p.log.Info(fmt.Sprintf("Start parsing. Find %v coins", len(crs)))

	p.initWorkers(taskChan)

	go func() {
		ticker := time.NewTicker(p.updateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.coins.Range(func(key, value any) bool {
					select {
					case taskChan <- value.(entity.Cryptocurrency):
						return true
					case <-p.done:
						return false
					}
				})
			case <-p.done:
				close(taskChan)
				return
			}
		}
	}()
}

func (p *Parser) initWorkers(taskChan chan entity.Cryptocurrency) {
	const op = "Parser.Workers"
	log := p.log.With(slog.String("op", op))

	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for coin := range taskChan {
				newPrice, err := p.cryptoClient.GetPrice(coin.Symbol, "USDT")
				if err != nil {
					log.Error(fmt.Sprintf("[Worker %d] Error fetching %s: %v\n", workerID, coin.Symbol, err))
					continue
				}

				ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
				defer done()
				hist := &entity.PriceHistory{
					CryptocurrencyID: coin.ID,
					Price:            newPrice,
					Timestamp:        time.Now(),
				}
				_, err = p.hst.Create(ctx, hist)
				if err != nil {
					log.Error(fmt.Sprintf("[Worker %d] Error creating price history %s: %v\n", workerID, coin.Symbol, err))
					continue
				}

				log.Info(fmt.Sprintf("[Worker %d] Updated %s: %.2f\n", workerID, coin.Symbol, newPrice))
			}
		}(i)
	}
}

func (p *Parser) Stop() {
	p.log.Info("Stop parsing")
	close(p.done)
	p.wg.Wait()
}

func (p *Parser) AddCoin(c entity.Cryptocurrency) {
	const op = "Parser.AddCoin"
	log := p.log.With(slog.String("op", op))
	p.coins.Store(c.Symbol, c)
	log.Info(fmt.Sprintf("Coin %s added to parser", c.Symbol))
}

func (p *Parser) RemoveCoin(c entity.Cryptocurrency) {
	const op = "Parser.RemoveCoin"
	log := p.log.With(slog.String("op", op))
	p.coins.Delete(c.Symbol)
	log.Info(fmt.Sprintf("Coin %s removed from parser", c.Symbol))
}
