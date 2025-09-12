package walletsrvc

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *WalletSuite) TestWithdrawService() {
	ctx := context.Background()
	userID := model.UserID(rand.Int64())
	available := rand.Int64N(math.MaxInt32)
	total := available + rand.Int64N(math.MaxInt32)
	amount := rand.Int64N(math.MaxInt16)
	description := fmt.Sprintf("desc_%d", rand.Int())
	txID := model.TxID(rand.Int64N(math.MaxInt16))
	refID := uuid.New()

	s.generatorService.EXPECT().NewUUID().Return(refID)
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
	).Return(nil)

	got, err := s.service.WithdrawalService(ctx, userID, amount, &description)
	s.NoError(err)
	s.Equal(refID, *got)
}
