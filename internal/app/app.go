package app

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Homyakadze14/AFFARM_tz/internal/config"
	v1 "github.com/Homyakadze14/AFFARM_tz/internal/controller/rest/v1"
	"github.com/Homyakadze14/AFFARM_tz/internal/infra/background"
	"github.com/Homyakadze14/AFFARM_tz/internal/infra/http"
	psg "github.com/Homyakadze14/AFFARM_tz/internal/infra/postgres"
	services "github.com/Homyakadze14/AFFARM_tz/internal/usecase"
	"github.com/Homyakadze14/AFFARM_tz/pkg/httpserver"
	"github.com/Homyakadze14/AFFARM_tz/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	s   *httpserver.Server
	p   *background.Parser
	db  *postgres.Postgres
	log *slog.Logger
}

func Run(
	log *slog.Logger,
	cfg *config.Config,
) *HttpServer {
	// Database
	pg, err := postgres.New(cfg.Database.URL, postgres.MaxPoolSize(cfg.Database.PoolMax))
	if err != nil {
		log.Error(fmt.Errorf("app - Run - postgres.New: %w", err).Error())
		os.Exit(1)
	}

	// Repository
	cryptocurRepo := psg.NewCryptocurrencyRepository(pg)
	trakingRepo := psg.NewTrackingRepository(pg)
	historyRepo := psg.NewHistoryRepository(pg)

	// Client
	timeout := 5 * time.Second
	binanceClient := http.NewBinanceClient(log, timeout)
	updateInterval := 5 * time.Second
	parser := background.NewParser(log, updateInterval, 10, historyRepo, cryptocurRepo, binanceClient)

	// Services
	cryptocurService := services.NewCryptocurrencyService(log, cryptocurRepo, trakingRepo, historyRepo, binanceClient, parser)

	// Parser
	go func() {
		parser.Start()
	}()

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(log, handler, cryptocurService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	return &HttpServer{s: httpServer, db: pg, log: log, p: parser}
}

func (s *HttpServer) Shutdown() {
	defer s.db.Close()
	defer s.p.Stop()
	err := s.s.Shutdown()
	if err != nil {
		s.log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}
