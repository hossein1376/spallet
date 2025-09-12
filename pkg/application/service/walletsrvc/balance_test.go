package walletsrvc

import (
	"context"
	"math"
	"math/rand/v2"

	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *WalletSuite) TestBalanceService() {
	ctx := context.Background()
	userID := model.UserID(rand.Int64())
	available := rand.Int64N(math.MaxInt16)
	total := available + rand.Int64N(math.MaxInt16)
	s.balanceRepo.EXPECT().
		Calculate(ctx, userID).
		Return(model.Balance{Available: available, Total: total}, nil).
		Once()

	got, err := s.service.BalanceService(ctx, userID)
	s.NoError(err)
	s.Equal(model.Balance{Available: available, Total: total}, got)
}
