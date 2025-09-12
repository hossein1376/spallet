package domain

import (
	"context"

	"github.com/google/uuid"
)

type Generator interface {
	NewUUID() uuid.UUID
}

type Worker interface {
	Run()
	Stop(ctx context.Context) error
	Add(
		ctx context.Context, id string, job, fallback func(context.Context) error,
	) error
}
