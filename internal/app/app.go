package app

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Homyakadze14/AFFARM_tz/internal/config"
	v1 "github.com/Homyakadze14/AFFARM_tz/internal/controller/rest/v1"
	"github.com/Homyakadze14/AFFARM_tz/internal/infra/http"
	psg "github.com/Homyakadze14/AFFARM_tz/internal/infra/postgres"
	services "github.com/Homyakadze14/AFFARM_tz/internal/usecase"
	"github.com/Homyakadze14/AFFARM_tz/pkg/httpserver"
	"github.com/Homyakadze14/AFFARM_tz/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	s   *httpserver.Server
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

	// Services
	cryptocurService := services.NewCryptocurrencyService(log, cryptocurRepo, trakingRepo, historyRepo, binanceClient)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, cryptocurService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	return &HttpServer{s: httpServer, db: pg, log: log}
}

func (s *HttpServer) Shutdown() {
	defer s.db.Close()
	err := s.s.Shutdown()
	if err != nil {
		s.log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}
