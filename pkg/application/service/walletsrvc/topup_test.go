package walletsrvc

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *WalletSuite) TestTopUpService() {
	ctx := context.Background()
	userID := model.UserID(rand.Int64())
	amount := rand.Int64N(math.MaxInt16)
	releaseDate := time.Now().Add(time.Second * time.Duration(rand.IntN(math.MaxInt16)))
	description := fmt.Sprintf("desc_%d", rand.Int())
	s.txRepo.EXPECT().
		Insert(
			ctx,
			userID,
			amount,
			model.TxTypeDeposit,
			model.InsertTxOption{
				ReleaseDate: &releaseDate, Description: &description},
		).
		Return(model.TxID(0), nil).
		Once()

	err := s.service.TopUpService(ctx, userID, amount, &releaseDate, &description)
	s.NoError(err)
}
