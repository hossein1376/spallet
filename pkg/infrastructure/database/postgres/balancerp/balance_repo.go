package balancerp

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/tools/errs"
)

var (
	ErrMissing = errors.New("user balance was not found")
)

const (
	sqlStateException = "P0001"
)

type BalanceRepo struct {
	tx *sql.Tx
}

func NewRepo(tx *sql.Tx) *BalanceRepo {
	return &BalanceRepo{tx: tx}
}

func (r *BalanceRepo) InsertZero(ctx context.Context, userID model.UserID) error {
	_, err := r.tx.ExecContext(
		ctx, "INSERT INTO balances (user_id) VALUES ($1);", userID,
	)
	return err
}

func (r *BalanceRepo) Calculate(
	ctx context.Context, userID model.UserID,
) (model.Balance, error) {
	var balance model.Balance
	err := r.tx.QueryRowContext(
		ctx,
		"SELECT total_balance, available_balance FROM refresh_user_balance($1);",
		userID,
	).Scan(&balance.Total, &balance.Available)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == sqlStateException {
			return balance, errs.NotFound(errs.WithErr(ErrMissing))
		}
		return balance, err
	}
	return balance, nil
}
