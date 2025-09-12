package workers

import (
	"context"
	"fmt"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/handler/config"
	"github.com/hossein1376/spallet/pkg/tools/worker"
)

var (
	_ domain.Worker = &worker.Worker{}
)

type Workers struct {
	Withdraw domain.Worker
}

func NewWorkers(cfg config.Worker) (*Workers, error) {
	withdraw, err := withdrawWorker(cfg)
	if err != nil {
		return nil, fmt.Errorf("new withdraw worker: %w", err)
	}

	return &Workers{Withdraw: withdraw}, nil
}

func (w *Workers) Run() {
	go func() {
		w.Withdraw.Run()
	}()
}

func (w *Workers) Stop(ctx context.Context) error {
	return w.Withdraw.Stop(ctx)
}
