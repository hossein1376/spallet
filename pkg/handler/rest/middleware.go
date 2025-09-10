package rest

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/hossein1376/spallet/pkg/tools/reqid"
	"github.com/hossein1376/spallet/pkg/tools/slogger"
)

func withMiddlewares(
	handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	var h http.Handler = handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := reqid.NewRequestID()
		ctx := context.WithValue(r.Context(), reqid.RequestIDKey, id)
		ctx = slogger.WithAttrs(ctx, slog.Any("request_id", id))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if msg := recover(); msg != nil {
				slogger.Error(
					r.Context(),
					"recovered panic",
					slog.Any("message", msg),
					slog.String("stack_trace", string(debug.Stack())),
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		raw := r.URL.RawQuery
		writer := &customWriter{ResponseWriter: w}
		defer func() {
			if raw != "" {
				path = path + "?" + raw
			}
			slogger.Info(r.Context(), "http request",
				slog.Group(
					"request",
					slog.String("method", r.Method),
					slog.String("path", path),
					slog.String("remote_addr", r.RemoteAddr), // In production, get from header
					slog.String("user_agent", r.UserAgent()),
				),
				slog.Group(
					"response",
					slog.Int("status", writer.status),
					slog.String("elapsed", time.Since(start).String()),
				),
			)
		}()

		next.ServeHTTP(writer, r)
	})
}

type customWriter struct {
	http.ResponseWriter
	status int
}

func (w *customWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}
