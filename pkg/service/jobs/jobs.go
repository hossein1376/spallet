package jobs

import (
	"fmt"

	"github.com/hossein1376/spallet/pkg/handler/config"
	"github.com/hossein1376/spallet/pkg/tools/worker"
)

type Jobs struct {
	WalletWithdraw *worker.Worker
}

func NewJobs(cfg config.Worker) (*Jobs, error) {
	w, err := worker.NewWorker(
		worker.WithRetryCount(cfg.RetryCount),
		worker.WithBufferSize(cfg.Size),
		worker.WithRetryInterval(cfg.DelayInterval),
		worker.WithBackoffMultiplier(cfg.BackoffMultiplier),
	)
	if err != nil {
		return nil, fmt.Errorf("wallet_withdraw worker: %w", err)
	}
	return &Jobs{WalletWithdraw: w}, nil
}

func (w *Jobs) Run() {
	go func() {
		w.WalletWithdraw.Run()
	}()
}

func (w *Jobs) Stop() {
	w.WalletWithdraw.Stop()
}
