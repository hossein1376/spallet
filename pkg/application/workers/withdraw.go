package workers

import (
	"fmt"

	"github.com/hossein1376/spallet/pkg/handler/config"
	"github.com/hossein1376/spallet/pkg/tools/worker"
)

func withdrawWorker(cfg config.Worker) (*worker.Worker, error) {
	w, err := worker.NewWorker(
		worker.WithRetryCount(cfg.RetryCount),
		worker.WithBufferSize(cfg.Size),
		worker.WithRetryInterval(cfg.DelayInterval),
		worker.WithBackoffMultiplier(cfg.BackoffMultiplier),
	)
	if err != nil {
		return nil, fmt.Errorf("wallet_withdraw worker: %w", err)
	}
	return w, nil
}
