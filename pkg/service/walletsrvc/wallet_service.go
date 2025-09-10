package walletsrvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/infrastructure/gateway"
	"github.com/hossein1376/spallet/pkg/service/jobs"
)

var (
	ErrInsufficientFunds = errors.New("insufficient balance to withdraw")
)

type WalletsService struct {
	pool    domain.Pool
	gateway gateway.Gateway
	jobs    *jobs.Jobs
}

func New(pool domain.Pool, j *jobs.Jobs) (*WalletsService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := pool.Query(ctx, func(ctx context.Context, r *domain.Repository) error {
		return r.Tx.RefundPending(ctx, nil)
	})
	if err != nil {
		return nil, fmt.Errorf("refund pending: %w", err)
	}

	return &WalletsService{pool: pool, gateway: gateway.New(), jobs: j}, nil
}
