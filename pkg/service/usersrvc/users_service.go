package usersrvc

import (
	"context"
	"fmt"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/infrastructure/gateway"
)

type UsersService struct {
	pool    domain.Pool
	gateway gateway.Gateway
}

func New(pool domain.Pool) *UsersService {
	return &UsersService{pool: pool, gateway: gateway.New()}
}

func (s UsersService) CreateUserService(
	ctx context.Context, username string,
) (*model.User, error) {
	var user *model.User
	q := func(ctx context.Context, r *domain.Repository) error {
		var err error
		user, err = r.Users.Create(ctx, username)
		if err != nil {
			return fmt.Errorf("creating user: %w", err)
		}
		err = r.Balance.InsertZero(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("creating new balance row: %w", err)
		}
		return nil
	}
	err := s.pool.Query(ctx, q)
	return user, err
}
