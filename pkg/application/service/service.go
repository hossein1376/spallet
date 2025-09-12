package service

import (
	"fmt"

	"github.com/hossein1376/spallet/pkg/application/generator"
	"github.com/hossein1376/spallet/pkg/application/service/usersrvc"
	"github.com/hossein1376/spallet/pkg/application/service/walletsrvc"
	"github.com/hossein1376/spallet/pkg/application/workers"
	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/handler/config"
	"github.com/hossein1376/spallet/pkg/infrastructure/gateway"
)

type Services struct {
	Wallets *walletsrvc.WalletsService
	Users   *usersrvc.UsersService
}

func NewServices(pool domain.Pool, cfg config.Worker) (
	*Services, *workers.Workers, error,
) {
	wrkrs, err := workers.NewWorkers(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("new workers: %w", err)
	}

	wallet, err := walletsrvc.New(
		pool, wrkrs.Withdraw, gateway.New(), generator.New(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("wallet service: %w", err)
	}
	users := usersrvc.New(pool)

	return &Services{Wallets: wallet, Users: users}, wrkrs, nil
}
