package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"
)

var (
	ErrNilWork     = errors.New("given work is nil")
	ErrQueueClosed = errors.New("queue is closed")
)

type Worker struct {
	logger     *slog.Logger
	wg         sync.WaitGroup
	jobs       chan job
	bufferSize int
	isClosed   bool

	retryCount        int
	delayInterval     time.Duration
	backoffMultiplier time.Duration
}

func NewWorker(opts ...Options) (*Worker, error) {
	q := &Worker{
		retryCount:    4,
		bufferSize:    1024,
		logger:        slog.Default(),
		delayInterval: time.Second,
	}
	for _, opt := range opts {
		if err := opt(q); err != nil {
			return nil, fmt.Errorf("option: %w", err)
		}
	}

	q.jobs = make(chan job, q.bufferSize)
	return q, nil
}

type job struct {
	ctx      context.Context
	id       string
	work     func(ctx context.Context) error
	fallback func(ctx context.Context) error
}

func (q *Worker) Add(
	ctx context.Context, id string, work, fallback func(ctx context.Context) error,
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

func (q *Worker) Run() {
	for jb := range q.jobs {
		q.wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					q.logger.ErrorContext(
						jb.ctx,
						"job panic",
						slog.String("job_id", jb.id),
						slog.Any("message", r),
						slog.String("stack", string(debug.Stack())),
					)
				}
			}()
			ctx := jb.ctx
			for attempt := range q.retryCount {
				err := jb.work(ctx)
				if err == nil {
					q.logger.DebugContext(
						ctx,
						"job completed",
						slog.String("job_id", jb.id),
					)
					return
				}
				q.logger.WarnContext(
					ctx,
					"job attempt failure",
					slog.String("job_id", jb.id),
					slog.Int("attempt", attempt+1),
					slog.String("error", err.Error()),
				)
				// Sleep between attempts
				if attempt < q.retryCount-1 {
					time.Sleep(q.delayInterval * q.backoffMultiplier)
				}
			}
			q.logger.WarnContext(ctx, "job failed", slog.String("job_id", jb.id))
			if jb.fallback == nil {
				return
			}
			if err := jb.fallback(ctx); err != nil {
				q.logger.ErrorContext(
					ctx,
					"job fallback func failure",
					slog.String("job_id", jb.id),
					slog.String("error", err.Error()),
				)
			}
		})
	}
}

func (q *Worker) Stop(ctx context.Context) error {
	doneCh := make(chan struct{})
	go func() {
		q.isClosed = true
		q.wg.Wait()
		close(q.jobs)
		doneCh <- struct{}{}
	}()
	select {
	case <-doneCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
