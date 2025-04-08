package main

import (
	"context"
	"log/slog"

	"pvz-service/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		panic(err)
	}

	err = a.Run()
	if err != nil {
		slog.Error("failed to run app", "error", err)
		panic(err)
	}
}
