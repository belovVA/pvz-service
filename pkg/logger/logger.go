package logger

import (
	"context"
	"log/slog"
	"os"

	"pvz-service/internal/middleware"
)

const (
	ErrorKey string = "error"
)

type HandlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddlware {
	return &HandlerMiddlware{next: next}
}

func (h *HandlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	// Добавляем userID, если он есть
	if userID, ok := ctx.Value(middleware.UserIDKey).(string); ok && userID != "" {
		rec.Add(middleware.UserIDKey, slog.StringValue(userID))
	}

	// Добавляем роль, если она есть
	if role, ok := ctx.Value(middleware.RoleKey).(string); ok && role != "" {
		rec.Add(middleware.RoleKey, slog.StringValue(role))
	}

	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)}
}

func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)}
}

func InitLogger() *slog.Logger {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	wrappedHandler := NewHandlerMiddleware(baseHandler)
	logger := slog.New(wrappedHandler)
	slog.SetDefault(logger)

	return logger
}
