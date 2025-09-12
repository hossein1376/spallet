package integration

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func NewMockDB(
	ctx context.Context, t *testing.T,
) (db *sql.DB, cleanup func(), err error) {
	dfrs := make([]func() error, 0)
	cleanup = func() {
		for _, dfr := range slices.Backward(dfrs) {
			if err := dfr(); err != nil {
				t.Error(err)
				return
			}
		}
	}

	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:17.6-trixie",
		postgres.WithDatabase("wallet"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("123456"),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(err)
	}
	dfrs = append(dfrs, func() error {
		if err := postgresContainer.Terminate(ctx); err != nil {
			return fmt.Errorf("postgresContainer.Terminate: %w", err)
		}
		return nil
	})

	db, err = sql.Open("pgx", postgresContainer.MustConnectionString(ctx))
	if err != nil {
		err = fmt.Errorf("open db connection: %w", err)
		return
	}
	dfrs = append(dfrs, func() error {
		if err := db.Close(); err != nil {
			return fmt.Errorf("close db connection: %w", err)
		}
		return nil
	})

	err = db.Ping()
	if err != nil {
		err = fmt.Errorf("ping db: %w", err)
		return
	}

	return
}
