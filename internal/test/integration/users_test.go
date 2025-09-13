package integration

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
)

func (s *IntegrationSuite) TestUsersCreateUser() {
	username := fmt.Sprintf("user_%d", rand.Int())
	var createdUser *model.User
	err := s.pool.Query(s.ctx, func(ctx context.Context, r *domain.Repository) error {
		var err error
		createdUser, err = r.Users.Create(ctx, username)
		return err
	})
	s.NoError(err)
	s.NotZero(createdUser.ID)
	s.Equal(username, createdUser.Username)

	var dbUser model.User
	err = s.db.QueryRowContext(
		s.ctx,
		"SELECT id, username FROM users WHERE username = $1",
		username,
	).Scan(&dbUser.ID, &dbUser.Username)
	s.NoError(err)
	s.Equal(*createdUser, dbUser)
}
