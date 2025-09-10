package txrp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/tools/errs"
)

var (
	ErrMissingUser = errors.New("user not found")
)

const (
	sqlStateConstraint = "23503"
)

type TxRepo struct {
	tx *sql.Tx
}

func NewRepo(tx *sql.Tx) *TxRepo {
	return &TxRepo{tx: tx}
}

func (r *TxRepo) Insert(
	ctx context.Context,
	userID model.UserID,
	amount int64,
	txType model.TxType,
	opt model.InsertTxOption,
) (model.TxID, error) {
	query :=
		`INSERT INTO transactions (
        	user_id, amount, type, status, ref_id, release_date, description
        )
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	args := []any{
		userID, amount, txType, opt.Status, opt.RefID, opt.ReleaseDate, opt.Description,
	}
	var id model.TxID
	err := r.tx.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == sqlStateConstraint {
			return 0, errs.NotFound(errs.WithErr(ErrMissingUser))
		}
		return 0, fmt.Errorf("inserting row: %w", err)
	}

	return id, nil
}

func (r *TxRepo) SetStaus(
	ctx context.Context,
	id model.TxID,
	status model.TxStatus,
	now time.Time,
) error {
	query := `UPDATE transactions SET status = $1, updated_at = $2
WHERE id = $3 AND status = 'pending';`
	_, err := r.tx.ExecContext(ctx, query, status, now, id)
	return err
}

func (r *TxRepo) ForUpdate(
	ctx context.Context, id model.TxID,
) (model.Transaction, error) {
	query := `SELECT id, user_id, amount, type, status, release_date, description,
	ref_id, created_at, updated_at
FROM transactions WHERE id = $1 FOR UPDATE;`
	return parseTransaction(r.tx.QueryRowContext(ctx, query, id))
}

func (r *TxRepo) List(
	ctx context.Context, userID model.UserID, count, threshold int64,
) ([]model.Transaction, error) {
	if threshold == 0 {
		threshold = 1<<63 - 1 // int64 max value
	}
	transactions := make([]model.Transaction, 0, count)
	query := `SELECT id, user_id, amount, type, status, release_date,description,
	ref_id, created_at, updated_at
FROM transactions
WHERE user_id = $1 AND id < $2
ORDER BY created_at DESC
LIMIT $3;`
	rows, err := r.tx.QueryContext(ctx, query, userID, threshold, count)
	if err != nil {
		return nil, fmt.Errorf("querying transactions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tx, err := parseTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating transactions: %w", err)
	}
	return transactions, nil
}

func (r *TxRepo) RefundPending(ctx context.Context, id *model.TxID) error {
	_, err := r.tx.ExecContext(ctx, "CALL refund_pending_transactions($1);", id)
	return err
}

func parseTransaction[
	T interface{ Scan(...any) error },
](row T) (model.Transaction, error) {
	var tx model.Transaction
	err := row.Scan(
		&tx.ID,
		&tx.UserID,
		&tx.Amount,
		&tx.Type,
		&tx.Status,
		&tx.ReleaseDate,
		&tx.Description,
		&tx.RefID,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	return tx, err
}
