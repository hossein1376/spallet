package service

import (
	"fmt"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/service/jobs"
	"github.com/hossein1376/spallet/pkg/service/usersrvc"
	"github.com/hossein1376/spallet/pkg/service/walletsrvc"
)

type Services struct {
	Wallets *walletsrvc.WalletsService
	Users   *usersrvc.UsersService
}

func NewServices(pool domain.Pool, w *jobs.Jobs) (*Services, error) {
	wallet, err := walletsrvc.New(pool, w)
	if err != nil {
		return nil, fmt.Errorf("wallet service: %w", err)
	}
	users := usersrvc.New(pool)

	return &Services{Wallets: wallet, Users: users}, nil
}
