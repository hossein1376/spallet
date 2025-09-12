package integration

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/infrastructure/database/postgres"
)

type TransactionsIntegrationSuite struct {
	suite.Suite

	ctx  context.Context
	db   *sql.DB
	pool domain.Pool
}

func TestTransactionsIntegrationSuite(t *testing.T) {
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

	s := &TransactionsIntegrationSuite{ctx: ctx, db: db, pool: pool}
	suite.Run(t, s)
}

func (s *TransactionsIntegrationSuite) TestInsert() {
	user := s.createUser()
	amount := rand.Int64N(math.MaxInt16)
	description := fmt.Sprintf("desc_%d", rand.Uint64())
	releaseDate := time.Now().
		Add(time.Duration(rand.IntN(math.MaxInt16))).
		Truncate(time.Millisecond)

	var txID model.TxID
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		var err error
		txID, err = r.Tx.Insert(
			ctx,
			user.ID,
			amount,
			model.TxTypeDeposit,
			model.InsertTxOption{
				Description: &description,
				ReleaseDate: &releaseDate,
			},
		)
		return err
	})
	s.NoError(err)
	s.NotZero(txID)

	var actual = s.getTx(txID)
	s.Equal(amount, actual.Amount)
	s.Equal(model.TxTypeDeposit, actual.Type)
	s.Equal(description, *actual.Description)
	s.Equal(releaseDate, *actual.ReleaseDate)
	s.Nil(actual.RefID)
	s.Nil(actual.Status)
	s.NotEmpty(actual.CreatedAt)
	s.NotEmpty(actual.UpdatedAt)
}

func (s *TransactionsIntegrationSuite) createUser() *model.User {
	var user *model.User
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		var err error
		user, err = r.Users.Create(ctx, fmt.Sprintf("user_%d", rand.Uint64()))
		if err != nil {
			return fmt.Errorf("creating user: %v", err)
		}
		err = r.Balance.InsertZero(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("inserting into balance table: %v", err)
		}
		return nil
	})
	s.NoError(err)
	return user
}

func (s *TransactionsIntegrationSuite) getTx(id model.TxID) model.Transaction {
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
