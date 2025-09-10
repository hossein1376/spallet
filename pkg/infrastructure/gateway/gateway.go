package gateway

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

var ErrRemoteGateway = fmt.Errorf("remote payment gateway error")

type Gateway struct{}

func New() Gateway {
	return Gateway{}
}

func (Gateway) Process(ctx context.Context, refID string) error {
	errCh := make(chan error, 1)
	go func() {
		var err error

		// 25% chance of failure.
		r := rand.IntN(4)
		if r == 0 {
			err = ErrRemoteGateway
		}

		// Simulate processing, and wait 1 to 4 seconds
		_ = refID
		time.Sleep(time.Duration(r+1) * time.Second)

		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
