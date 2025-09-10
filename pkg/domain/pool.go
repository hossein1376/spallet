package domain

import (
	"context"
)

type Pool interface {
	Query(ctx context.Context, f QueryFunc) error
}

type QueryFunc = func(ctx context.Context, r *Repository) error
