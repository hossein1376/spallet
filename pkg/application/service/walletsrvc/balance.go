package walletsrvc

import (
	"context"
	"fmt"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *WalletsService) BalanceService(
	ctx context.Context, userID model.UserID,
) (
	model.Balance, error,
) {
	var balance model.Balance
	q := func(ctx context.Context, r *domain.Repository) error {
		var err error
		balance, err = r.Balance.Calculate(ctx, userID)
		return err
	}

	err := s.pool.Query(ctx, q)
	if err != nil {
		return balance, fmt.Errorf("fetching balance: %w", err)
	}
	return balance, nil
}
