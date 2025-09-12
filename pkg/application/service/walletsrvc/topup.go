package walletsrvc

import (
	"context"
	"fmt"
	"time"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *WalletsService) TopUpService(
	ctx context.Context,
	userID model.UserID,
	amount int64,
	releaseDate *time.Time,
	description *string,
) error {
	q := func(ctx context.Context, r *domain.Repository) error {
		_, err := r.Tx.Insert(
			ctx,
			userID,
			amount,
			model.TxTypeDeposit,
			model.InsertTxOption{
				ReleaseDate: releaseDate,
				Description: description,
			},
		)
		if err != nil {
			return fmt.Errorf("inserting wallet transaction: %w", err)
		}
		return nil
	}

	return s.pool.Query(ctx, q)
}
