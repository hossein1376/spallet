package walletsrvc

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/hossein1376/spallet/internal/mocks"
	"github.com/hossein1376/spallet/internal/tester"
)

type WalletSuite struct {
	suite.Suite
	service          *WalletsService
	txRepo           *mocks.MockTransactionsRepository
	balanceRepo      *mocks.MockBalanceRepository
	usersRepo        *mocks.MockUsersRepository
	generatorService *mocks.MockGenerator
	worker           *mocks.MockWorker
}

func TestWalletSuite(t *testing.T) {
	suite.Run(t, new(WalletSuite))
}

func (s *WalletSuite) SetupSuite() {
	repo := tester.NewMockRepo(s.T())
	pool := tester.NewMockPool(s.T(), repo)
	gateway := mocks.NewMockGateway(s.T())
	worker := mocks.NewMockWorker(s.T())
	generator := mocks.NewMockGenerator(s.T())

	s.txRepo = repo.Tx.(*mocks.MockTransactionsRepository)
	s.usersRepo = repo.Users.(*mocks.MockUsersRepository)
	s.balanceRepo = repo.Balance.(*mocks.MockBalanceRepository)

	s.txRepo.EXPECT().
		RefundPending(mock.Anything, mock.Anything).
		Return(nil).
		Once()

	srvc, err := New(pool, worker, gateway, generator)
	s.Require().NoError(err)
	s.NotNil(srvc)

	s.service = srvc
	s.worker = worker
	s.generatorService = generator
}
