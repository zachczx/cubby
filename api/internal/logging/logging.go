package logging

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

var RequestIDKey contextKey = "requestID"

func Init() {
	env := os.Getenv("ENV")
	var handler slog.Handler

	if env == "development" {
		handler = slog.NewTextHandler(os.Stdout, nil)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(handler))
}

func WithRequestID(ctx context.Context) *slog.Logger {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return slog.With("requestId", id)
	}
	return slog.Default()
}

func Info(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Info(msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Error(msg, args...)
}
