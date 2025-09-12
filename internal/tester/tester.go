package tester

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/hossein1376/spallet/internal/mocks"
	"github.com/hossein1376/spallet/pkg/domain"
)

func NewMockPool(t *testing.T, repo *domain.Repository) *mocks.MockPool {
	pool := mocks.NewMockPool(t)
	pool.EXPECT().
		Query(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, f domain.QueryFunc) error {
			return f(ctx, repo)
		})
	return pool
}

func NewMockRepo(t *testing.T) *domain.Repository {
	return &domain.Repository{
		Tx:      mocks.NewMockTransactionsRepository(t),
		Users:   mocks.NewMockUsersRepository(t),
		Balance: mocks.NewMockBalanceRepository(t),
	}
}
