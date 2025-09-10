package serde

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/hossein1376/spallet/pkg/tools/errs"
	"github.com/hossein1376/spallet/pkg/tools/slogger"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func ValueOrDefault[T any](v string, f func(string) (T, error)) (T, error) {
	if v == "" {
		var t T
		return t, nil
	}
	return f(v)
}

func ExtractFromErr(ctx context.Context, err error) (int, Response) {
	if err == nil {
		panic("ExtractFromErr was called with nil error")
	}

	var e errs.Error
	if errors.As(err, &e) {
		slogger.Debug(
			ctx,
			"Error response",
			slog.Int("status_code", e.HTTPStatusCode),
			slog.String("message", e.Message),
			slogger.Err("error", e.Err),
		)
		return e.HTTPStatusCode, Response{Message: e.Message}
	}

	slogger.Error(ctx, "Internal server error", slogger.Err("error", err))
	return http.StatusInternalServerError, Response{
		Message: http.StatusText(http.StatusInternalServerError),
	}
}
