package command

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/hossein1376/spallet/pkg/handler/config"
	"github.com/hossein1376/spallet/pkg/handler/rest"
	"github.com/hossein1376/spallet/pkg/infrastructure/database/postgres"
	"github.com/hossein1376/spallet/pkg/service"
	"github.com/hossein1376/spallet/pkg/service/jobs"
	"github.com/hossein1376/spallet/pkg/tools/slogger"
)

func Run() error {
	ctx := context.Background()

	var cfgPath string
	flag.StringVar(&cfgPath, "config", "assets/config.yaml", "config file path")
	flag.Parse()

	logger := slogger.NewJSONLogger(slog.LevelDebug, os.Stdout)
	slog.SetDefault(logger)

	cfg, err := config.New(cfgPath)
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}

	pool, err := postgres.NewPool(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("new pool: %w", err)
	}
	defer pool.Close()

	w, err := jobs.NewJobs(cfg.Worker)
	if err != nil {
		return fmt.Errorf("new wallet_withdraw jobs: %w", err)
	}
	defer w.Stop()
	w.Run()

	services, err := service.NewServices(pool, w)
	if err != nil {
		return fmt.Errorf("new services: %w", err)
	}
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
