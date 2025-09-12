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

type BalancesIntegrationSuite struct {
	suite.Suite

	ctx  context.Context
	db   *sql.DB
	pool domain.Pool
}

func TestBalancesIntegrationSuite(t *testing.T) {
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

	s := &BalancesIntegrationSuite{ctx: ctx, db: db, pool: pool}
	suite.Run(t, s)
}

func (s *BalancesIntegrationSuite) TestInsertZero() {
	user := s.createUser(false)

	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		return r.Balance.InsertZero(ctx, user.ID)
	})
	s.NoError(err)

	var total, available int64
	var lastTime time.Time
	err = s.db.QueryRowContext(
		s.ctx,
		"SELECT total, available, last_calculated_at FROM balances WHERE user_id = $1",
		user.ID,
	).Scan(&total, &available, &lastTime)
	s.NoError(err)
	s.Zero(available)
	s.Zero(total)
	s.Equal(time.Unix(0, 0).Local(), lastTime)
}

func (s *BalancesIntegrationSuite) TestCalculate() {
	user := s.createUser(true)
	amount := rand.Int64N(math.MaxInt16)
	s.insertTx(user.ID, amount, model.TxTypeDeposit, model.InsertTxOption{})

	var calculated model.Balance
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		var err error
		calculated, err = r.Balance.Calculate(ctx, user.ID)
		return err
	})
	s.NoError(err)
	s.NotZero(calculated)
	s.Equal(amount, calculated.Total)
	s.Equal(amount, calculated.Available)

	var dbSum int64
	err = s.db.QueryRowContext(
		s.ctx,
		"SELECT SUM(CASE WHEN type='deposit' THEN amount ELSE -amount END) FROM transactions WHERE user_id = $1",
		user.ID,
	).Scan(&dbSum)
	s.NoError(err)
}

func (s *BalancesIntegrationSuite) createUser(createBalance bool) *model.User {
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

func (s *BalancesIntegrationSuite) insertTx(
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
