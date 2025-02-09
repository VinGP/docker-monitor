package app

import (
	"backend/internal/config"
	"backend/internal/http/handlers"
	"backend/internal/kafka/consumer"
	"backend/internal/repo"
	"backend/internal/service"
	"backend/pkg/httpserver"
	"backend/pkg/logger"
	"backend/pkg/logger/sl"
	"backend/pkg/postgres"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Config) {
	log := logger.NewLogger(cfg.Production)

	log.Info("app config", slog.Any("config", cfg))

	log.Info(
		"starting user-api",
		slog.Bool("PRODUCTION", cfg.Production),
	)
	log.Debug("debug messages are enabled")

	r := chi.NewRouter()

	postgresDB, err := postgres.New(cfg.DBUrl)
	if err != nil {
		log.Error("postgres.New", sl.Err(err))
		os.Exit(1)
	}

	containerStatusRepo := repo.New(postgresDB)

	handlers.NewRouter(r, log, service.NewContainerStatusService(containerStatusRepo))

	httpServer := httpserver.New(r, httpserver.Port(cfg.Port))

	containerStatusConsumer, err := consumer.NewContainerStatusConsumer(cfg.ContainerStatusTopic,
		cfg.ContainerStatusGroupID,
		service.NewContainerStatusService(containerStatusRepo), cfg.Brokers)
	containerStatusConsumer.Start()

	if err != nil {
		log.Error("consumer.NewContainerStatusConsumer", sl.Err(err))
		os.Exit(1)
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: ", slog.Any("signal", s.String()))
	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify: %w", sl.Err(err))
	}

	if err = httpServer.Shutdown(); err != nil {
		log.Error("app - Run - httpServer.Shutdown: %w", sl.Err(err))
	}

	if err = containerStatusConsumer.Stop(); err != nil {
		log.Error("app - Run - containerStatusConsumer.Stop: %w", sl.Err(err))
	}

	postgresDB.Close()

	log.Info("app - Run - db closed")

	log.Info("app - Run - exiting")
}
