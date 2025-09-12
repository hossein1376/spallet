package integration

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite

	ctx context.Context
	db  *sql.DB
}

func TestIntegrationSuite(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := NewMockDB(ctx, t)
	defer cleanup()
	if err != nil {
		t.Errorf("creating DB container: %v", err)
		return
	}

	s := &IntegrationSuite{ctx: ctx, db: db}
	suite.Run(t, s)
}

type EndToEndSuite struct {
	suite.Suite

	// ...
}

func TestEndToEndSuite(t *testing.T) {
	// ...

	s := &EndToEndSuite{}
	suite.Run(t, s)
}
