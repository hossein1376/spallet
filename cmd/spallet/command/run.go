package command

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/hossein1376/spallet/pkg/application/service"
	"github.com/hossein1376/spallet/pkg/handler/config"
	"github.com/hossein1376/spallet/pkg/handler/rest"
	"github.com/hossein1376/spallet/pkg/infrastructure/database/postgres"
	"github.com/hossein1376/spallet/pkg/tools/slogger"
)

func Run() error {
	ctx := context.Background()

	var cfgPath string
	flag.StringVar(&cfgPath, "config", "assets/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.New(cfgPath)
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}

	logger := slogger.NewJSONLogger(cfg.Server.LogLevel, os.Stdout)
	slog.SetDefault(logger)

	pool, err := postgres.NewPool(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("new pool: %w", err)
	}
	defer func() {
		if err := pool.Close(); err != nil {
			slog.Warn("closing pool", slog.String("error", err.Error()))
		}
	}()

	services, worker, err := service.NewServices(pool, cfg.Worker)
	if err != nil {
		return fmt.Errorf("new services: %w", err)
	}
	defer func() {
		if err := worker.Stop(ctx); err != nil {
			slog.Warn("stopping worker", slog.String("error", err.Error()))
		}
	}()
	worker.Run()

	server := rest.NewServer(cfg.Server.Addr, services)
	slog.Debug("initialized repositories, services and handlers")

	errCh := make(chan error)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		slog.Info("starting server", slog.String("addr", cfg.Server.Addr))
		errCh <- server.ListenAndServe()
	}()

	select {
	case err = <-errCh:
		return fmt.Errorf("srv.ListenAndServe: %w", err)
	case <-signalCh:
		slog.Info("shutdown signal received")
		return server.Shutdown(ctx)
	}
}
