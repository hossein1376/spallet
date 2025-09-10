package domain

import (
	"context"
	"time"

	"github.com/hossein1376/spallet/pkg/domain/model"
)

type Repository struct {
	Tx      TransactionsRepository
	Users   UsersRepository
	Balance BalanceRepository
}

type TransactionsRepository interface {
	Insert(
		ctx context.Context,
		userID model.UserID,
		amount int64,
		txType model.TxType,
		opts model.InsertTxOption,
	) (model.TxID, error)
	ForUpdate(ctx context.Context, id model.TxID) (model.Transaction, error)
	SetStaus(
		ctx context.Context,
		id model.TxID,
		status model.TxStatus,
		now time.Time,
	) error
	List(
		ctx context.Context, userID model.UserID, count, threshold int64,
	) ([]model.Transaction, error)
	RefundPending(ctx context.Context, id *model.TxID) error
}

type BalanceRepository interface {
	InsertZero(ctx context.Context, userID model.UserID) error
	Calculate(ctx context.Context, userID model.UserID) (model.Balance, error)
}

type UsersRepository interface {
	Create(ctx context.Context, username string) (*model.User, error)
}
