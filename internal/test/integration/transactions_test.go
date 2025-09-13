package integration

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *IntegrationSuite) TestTransactions_Insert() {
	user := s.createUser(true)
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

func (s *IntegrationSuite) TestTransactions_List() {
	user := s.createUser(true)
	amountTx1 := rand.Int64N(math.MaxInt16)
	amountTx2 := rand.Int64N(math.MaxInt16)
	tx1 := s.insertTx(
		user.ID, amountTx1, model.TxTypeDeposit, model.InsertTxOption{},
	)
	tx2 := s.insertTx(
		user.ID,
		amountTx2,
		model.TxTypeWithdrawal,
		model.InsertTxOption{
			Status: model.Ptr(model.TxStatusPending),
		},
	)
	var list []model.Transaction
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		var err error
		list, err = r.Tx.List(ctx, user.ID, 10, 0)
		return err
	})
	s.NoError(err)
	s.Len(list, 2)

	s.Equal(tx2, list[0].ID)
	s.Equal(user.ID, list[0].UserID)
	s.Equal(amountTx2, list[0].Amount)
	s.Equal(model.TxStatusPending, *list[0].Status)
	s.Equal(model.TxTypeWithdrawal, list[0].Type)

	s.Equal(tx1, list[1].ID)
	s.Equal(user.ID, list[1].UserID)
	s.Equal(amountTx1, list[1].Amount)
	s.Equal(model.TxTypeDeposit, list[1].Type)
	s.Nil(list[1].Status)
}

func (s *IntegrationSuite) TestTransactions_SetStaus() {
	user := s.createUser(true)
	txID := s.insertTx(
		user.ID, rand.Int64(), model.TxTypeWithdrawal, model.InsertTxOption{
			Status: model.Ptr(model.TxStatusPending),
		},
	)
	now := time.Now().Truncate(time.Millisecond)
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		return r.Tx.SetStaus(ctx, txID, model.TxStatusCompleted, now)
	})
	s.NoError(err)

	actual := s.getTx(txID)
	s.Equal(model.TxStatusCompleted, *actual.Status)
	s.Equal(now, actual.UpdatedAt)
}

func (s *IntegrationSuite) TestTransactions_RefundPending() {
	user := s.createUser(true)
	amount := rand.Int64N(math.MaxInt16)
	refID := uuid.New()
	txID := s.insertTx(
		user.ID, amount, model.TxTypeWithdrawal, model.InsertTxOption{
			Status: model.Ptr(model.TxStatusPending),
			RefID:  &refID,
		},
	)
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		return r.Tx.RefundPending(ctx, &txID)
	})
	s.NoError(err)

	actual := s.getTx(txID)
	s.Equal(model.TxStatusFailed, *actual.Status)
	s.Equal(refID, *actual.RefID)

	refunded := s.getTx(txID + 1) // not a good idea, but it works for now
	s.Equal(model.TxTypeDeposit, refunded.Type)
	s.Equal(amount, refunded.Amount)
	s.Equal(fmt.Sprintf("refund %s", refID), *refunded.Description)
}
