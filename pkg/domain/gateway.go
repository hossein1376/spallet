package domain

import (
	"context"
)

type Gateway interface {
	Process(ctx context.Context, refID string) error
}
