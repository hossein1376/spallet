package walletsrvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hossein1376/spallet/pkg/domain"
)

var (
	ErrInsufficientFunds = errors.New("insufficient balance to withdraw")
)

type WalletsService struct {
	pool        domain.Pool
	gateway     domain.Gateway
	withdrawJob domain.Worker
	generator   domain.Generator
}

func New(
	pool domain.Pool,
	worker domain.Worker,
	gateway domain.Gateway,
	generator domain.Generator,
) (*WalletsService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := pool.Query(ctx, func(ctx context.Context, r *domain.Repository) error {
		return r.Tx.RefundPending(ctx, nil)
	})
	if err != nil {
		return nil, fmt.Errorf("refund pending: %w", err)
	}

	return &WalletsService{
		pool:        pool,
		gateway:     gateway,
		withdrawJob: worker,
		generator:   generator,
	}, nil
}
