package walletsrvc

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/tools/errs"
)

func (s *WalletSuite) TestWithdrawService() {
	ctx := context.Background()
	userID := model.UserID(rand.Int64())
	description := fmt.Sprintf("desc_%d", rand.Int())
	txID := model.TxID(rand.Int64N(math.MaxInt16))
	refID := uuid.New()

	s.Run("withdraw successfully", func() {
		available := rand.Int64N(math.MaxInt32)
		total := available + rand.Int64N(math.MaxInt32)
		amount := rand.Int64N(math.MaxInt16)

		s.generatorService.EXPECT().NewUUID().Return(refID).Once()
		s.balanceRepo.EXPECT().
			Calculate(ctx, userID).
			Return(model.Balance{Available: available, Total: total}, nil).
			Once()
		s.txRepo.EXPECT().
			Insert(
				ctx,
				userID,
				amount,
				model.TxTypeWithdrawal,
				model.InsertTxOption{
					Description: &description,
					RefID:       &refID,
					Status:      model.Ptr(model.TxStatusPending),
				},
			).
			Return(txID, nil).
			Once()
		s.worker.EXPECT().Add(
			context.WithoutCancel(ctx),
			refID.String(),
			mock.AnythingOfType("func(context.Context) error"),
			mock.AnythingOfType("func(context.Context) error"),
		).Return(nil).Once()

		got, err := s.service.WithdrawalService(ctx, userID, amount, &description)
		s.NoError(err)
		s.Equal(refID, *got)
	})
	s.Run("Not enough funds", func() {
		amount := rand.Int64N(math.MaxInt16) + 1
		s.balanceRepo.EXPECT().
			Calculate(ctx, userID).
			Return(model.Balance{Available: 0, Total: 0}, nil).
			Once()

		got, err := s.service.WithdrawalService(ctx, userID, amount, &description)
		s.Nil(got)
		s.ErrorIs(err, ErrInsufficientFunds)
		var conflictErr errs.Error
		s.ErrorAs(err, &conflictErr)
		s.Equal(http.StatusConflict, conflictErr.HTTPStatusCode)
	})
}
