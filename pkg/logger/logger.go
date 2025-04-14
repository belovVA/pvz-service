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

type Handler struct {
	next slog.Handler
}

func NewHandlerLogger(next slog.Handler) *Handler {
	return &Handler{next: next}
}

func (h *Handler) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
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

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{next: h.next.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{next: h.next.WithGroup(name)}
}

func InitLogger() *slog.Logger {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	wrappedHandler := NewHandlerLogger(baseHandler)
	logger := slog.New(wrappedHandler)
	slog.SetDefault(logger)

	return logger
}
