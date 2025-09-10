package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/hossein1376/spallet/pkg/domain"
	"github.com/hossein1376/spallet/pkg/handler/config"
	"github.com/hossein1376/spallet/pkg/infrastructure/database/postgres/balancerp"
	"github.com/hossein1376/spallet/pkg/infrastructure/database/postgres/txrp"
	"github.com/hossein1376/spallet/pkg/infrastructure/database/postgres/usersrp"
)

func newRepo(tx *sql.Tx) *domain.Repository {
	return &domain.Repository{
		Tx:      txrp.NewRepo(tx),
		Users:   usersrp.NewRepo(tx),
		Balance: balancerp.NewRepo(tx),
	}
}

type Pool struct {
	db *sql.DB
}

func NewPool(ctx context.Context, cfg config.DB) (*Pool, error) {
	tls := "enable"
	if cfg.DisableTLS {
		tls = "disable"
	}
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s database=%s sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		tls,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open(): %w", err)
	}

	// Check if the connection is established
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	pool := &Pool{db: db}
	if err = pool.migrate(ctx); err != nil {
		return nil, fmt.Errorf("migrating database: %w", err)
	}

	return pool, nil
}

func (p *Pool) Query(ctx context.Context, f domain.QueryFunc) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	err = f(ctx, newRepo(tx))
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rollback: %w query: %w", rollbackErr, err)
		}
		return fmt.Errorf("query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (p *Pool) Close() error {
	return p.db.Close()
}

//go:embed schema/schema.sql
var schema string

func (p *Pool) migrate(ctx context.Context) error {
	existenceCheck := "SELECT 0 FROM transactions LIMIT 1;"
	_, err := p.db.ExecContext(ctx, existenceCheck)
	if err == nil {
		return nil
	}
	_, err = p.db.ExecContext(ctx, schema)
	return err
}
