package worker

import (
	"fmt"
	"log/slog"
	"time"
)

type Options func(*Worker) error

func WithLogger(logger *slog.Logger) Options {
	return func(q *Worker) error {
		q.logger = logger
		return nil
	}
}

func WithBackoffMultiplier(multiplier int) Options {
	return func(q *Worker) error {
		if multiplier < 1 {
			return fmt.Errorf(
				"multiplier must be greater than equal 1: %d", multiplier,
			)
		}
		q.backoffMultiplier = time.Duration(multiplier)
		return nil
	}
}

func WithBufferSize(bufferSize int) Options {
	return func(q *Worker) error {
		if bufferSize <= 0 {
			return fmt.Errorf("buffer size must be positive: %d", bufferSize)
		}
		q.bufferSize = bufferSize
		return nil
	}
}

func WithRetryCount(retryCount int) Options {
	return func(q *Worker) error {
		if retryCount < 0 {
			return fmt.Errorf("retry count must be positive: %d", retryCount)
		}
		q.retryCount = retryCount
		return nil
	}
}

func WithRetryInterval(interval time.Duration) Options {
	return func(q *Worker) error {
		if interval < 0 {
			return fmt.Errorf("interval must be greater than 0: %d", interval)
		}
		q.delayInterval = interval
		return nil
	}
}
