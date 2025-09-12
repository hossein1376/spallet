package walletsrvc

import (
	"context"
	"math"
	"math/rand/v2"

	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *WalletSuite) TestHistoryService() {
	ctx := context.Background()
	userID := model.UserID(rand.Int64())
	count := rand.Int64N(math.MaxInt16)
	threshold := rand.Int64N(math.MaxInt16)
	s.txRepo.EXPECT().
		List(ctx, userID, count, threshold).
		Return([]model.Transaction{}, nil).
		Once()

	got, err := s.service.HistoryService(ctx, userID, count, threshold)
	s.NoError(err)
	s.NotNil(got)
}
