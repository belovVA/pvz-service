package pkg

import (
	"context"
	"log/slog"
	"os"
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
	//if c, ok := ctx.Value(key).(logCtx); ok {
	//	if c.UserID != 0 {
	//		rec.Add("userID", c.UserID)
	//	}
	//	if c.Phone != "" {
	//		rec.Add("phone", c.Phone)
	//	}
	//	if c.Gate != "" {
	//		rec.Add("sms_gate", c.Gate)
	//	}
	//	if c.Message != "" {
	//		rec.Add("message", c.Message)
	//	}
	//}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}
func InitLogger() {
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}
