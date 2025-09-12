package usersrvc_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hossein1376/spallet/internal/mocks"
	"github.com/hossein1376/spallet/internal/tester"
	"github.com/hossein1376/spallet/pkg/application/service/usersrvc"
	"github.com/hossein1376/spallet/pkg/domain/model"
)

func TestUsersService_CreateUserService(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	id := model.UserID(rand.Int64())
	username := fmt.Sprintf("user-%d", rand.Int())

	repo := tester.NewMockRepo(t)
	repo.Users.(*mocks.MockUsersRepository).
		EXPECT().
		Create(ctx, username).
		Return(&model.User{ID: id, Username: username}, nil).
		Once()
	repo.Balance.(*mocks.MockBalanceRepository).
		EXPECT().
		InsertZero(ctx, id).
		Return(nil).
		Once()

	service := usersrvc.New(tester.NewMockPool(t, repo))
	got, err := service.CreateUserService(ctx, username)
	a.NoError(err)
	a.NotNil(got)
	a.Equal(id, got.ID, "id")
	a.Equal(username, got.Username, "username")
}
