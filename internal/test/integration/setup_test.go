package integration

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/infrastructure/database/postgres"
)

type IntegrationSuite struct {
	suite.Suite

	ctx  context.Context
	db   *sql.DB
	pool domain.Pool
}

func TestIntegrationSuite(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := NewMockDB(ctx, t)
	defer cleanup()
	if err != nil {
		t.Errorf("creating DB container: %v", err)
		return
	}

	pool, err := postgres.NewFromDB(ctx, db)
	if err != nil {
		t.Errorf("creating pool: %v", err)
		return
	}

	s := &IntegrationSuite{ctx: ctx, db: db, pool: pool}
	suite.Run(t, s)
}

func (s *IntegrationSuite) createUser(createBalance bool) *model.User {
	var user *model.User
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		var err error
		user, err = r.Users.Create(ctx, fmt.Sprintf("user_%d", rand.Uint64()))
		if err != nil {
			return fmt.Errorf("creating user: %v", err)
		}
		if createBalance {
			err = r.Balance.InsertZero(ctx, user.ID)
			if err != nil {
				return fmt.Errorf("inserting into balance table: %v", err)
			}
		}
		return nil
	})
	s.NoError(err)
	return user
}

func (s *IntegrationSuite) insertTx(
	userID model.UserID,
	amount int64,
	txType model.TxType,
	opts model.InsertTxOption,
) model.TxID {
	var txID model.TxID
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		var err error
		txID, err = r.Tx.Insert(ctx, userID, amount, txType, opts)
		return err
	})
	s.NoError(err)
	return txID
}

func (s *IntegrationSuite) getTx(id model.TxID) model.Transaction {
	var tx model.Transaction
	err := s.db.QueryRowContext(
		s.ctx,
		`SELECT id, user_id, amount, type, status, release_date,description,
       		ref_id, created_at, updated_at FROM transactions WHERE id = $1`,
		id,
	).Scan(
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
	s.NoError(err)
	return tx
}
