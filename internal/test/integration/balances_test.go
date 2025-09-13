package integration

import (
	"context"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *IntegrationSuite) TestBalances_InsertZero() {
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

func (s *IntegrationSuite) TestBalances_Calculate() {
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
		`SELECT SUM(CASE WHEN type='deposit' THEN amount ELSE -amount END)
FROM transactions WHERE user_id = $1`,
		user.ID,
	).Scan(&dbSum)
	s.NoError(err)
}
