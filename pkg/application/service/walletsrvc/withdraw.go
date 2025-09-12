package walletsrvc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/tools/errs"
)

func (s *WalletsService) WithdrawalService(
	ctx context.Context, userID model.UserID, amount int64, description *string,
) (*uuid.UUID, error) {
	var (
		refID *uuid.UUID
		txID  model.TxID
	)
	err := s.pool.Query(ctx, func(ctx context.Context, r *domain.Repository) error {
		balance, err := r.Balance.Calculate(ctx, userID)
		if err != nil {
			return fmt.Errorf("calculating balance: %w", err)
		}
		if balance.Available < amount {
			return errs.Conflict(errs.WithErr(ErrInsufficientFunds))
		}

		refID = model.Ptr(s.generator.NewUUID())
		txID, err = r.Tx.Insert(
			ctx,
			userID,
			amount,
			model.TxTypeWithdrawal,
			model.InsertTxOption{
				Description: description,
				Status:      model.Ptr(model.TxStatusPending),
				RefID:       refID,
			},
		)
		if err != nil {
			return fmt.Errorf("inserting withdraw transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	err = s.withdrawJob.Add(
		context.WithoutCancel(ctx),
		refID.String(),
		s.withdraw(txID),
		s.refund(txID),
	)
	if err != nil {
		return nil, fmt.Errorf("enqueuing process request: %w", err)
	}

	return refID, nil
}

func (s *WalletsService) withdraw(txID model.TxID) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		q := func(ctx context.Context, r *domain.Repository) error {
			now := time.Now()
			tx, err := r.Tx.ForUpdate(ctx, txID)
			if err != nil {
				return fmt.Errorf("finding transaction: %w", err)
			}
			if tx.Type != model.TxTypeWithdrawal {
				return fmt.Errorf("invalid transaction type: %s", tx.Type)
			}
			err = s.gateway.Process(ctx, tx.RefID.String())
			if err != nil {
				return fmt.Errorf("processing request: %w", err)
			}
			err = r.Tx.SetStaus(ctx, txID, model.TxStatusCompleted, now)
			if err != nil {
				return fmt.Errorf("setting transaction staus: %w", err)
			}
			return nil
		}
		return s.pool.Query(ctx, q)
	}
}

func (s *WalletsService) refund(txID model.TxID) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return s.pool.Query(
			ctx,
			func(ctx context.Context, r *domain.Repository) error {
				return r.Tx.RefundPending(ctx, &txID)
			},
		)
	}
}
