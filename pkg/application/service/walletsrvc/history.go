package walletsrvc

import (
	"context"
	"fmt"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *WalletsService) HistoryService(
	ctx context.Context, userID model.UserID, count, threshold int64,
) ([]model.Transaction, error) {
	var transactions []model.Transaction
	q := func(ctx context.Context, r *domain.Repository) error {
		var err error
		transactions, err = r.Tx.List(ctx, userID, count, threshold)
		return err
	}
	err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("fetching transactions: %w", err)
	}
	return transactions, nil
}
