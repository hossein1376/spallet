package integration

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/domain/model"
	"github.com/hossein1376/spallet/pkg/infrastructure/database/postgres"
)

type UsersIntegrationSuite struct {
	suite.Suite

	ctx  context.Context
	db   *sql.DB
	pool domain.Pool
}

func TestUsersIntegrationSuite(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := NewMockDB(ctx, t)
	defer cleanup()
	if err != nil {
		t.Errorf("creating DB container: %v", err)
		return
	}

	pool, err := postgres.NewFromDB(ctx, db)
	if err != nil {
		t.Errorf("creating pool: %v", err)
		return
	}

	s := &UsersIntegrationSuite{ctx: ctx, db: db, pool: pool}
	suite.Run(t, s)
}

func (s *UsersIntegrationSuite) TestCreateUser() {
	s.cleanUsersTable()

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

func (s *UsersIntegrationSuite) cleanUsersTable() {
	_, err := s.db.ExecContext(s.ctx, "DELETE FROM users;")
	s.NoError(err)
}
