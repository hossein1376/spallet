package worker

import (
	"context"
	"log/slog"
)

type job struct {
	ctx            context.Context
	id             string
	work, fallback Work
}

type Work func(ctx context.Context) error

func (q *Worker) Add(
	ctx context.Context, id string, work, fallback Work,
) error {
	switch {
	case q.isClosed:
		return ErrQueueClosed
	case work == nil:
		return ErrNilWork
	}

	go func() {
		q.jobs <- job{ctx: ctx, id: id, work: work, fallback: fallback}
	}()
	q.logger.DebugContext(ctx, "add new job", slog.String("id", id))
	return nil
}
