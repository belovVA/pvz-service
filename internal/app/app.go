package app

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pvz-service/internal/handler"
	"pvz-service/internal/middleware"
	"pvz-service/internal/repository"
	"pvz-service/internal/service"

	"github.com/go-chi/chi/v5"
	"pvz-service/internal/config"
	"pvz-service/pkg"
)

type App struct {
	httpCfg config.HTTPConfig
	router  *chi.Mux
}

func NewApp(ctx context.Context) (*App, error) {
	pkg.InitLogger()

	pgCfg, err := config.PGConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading postgres config: %w", err)
	}

	htppCfg, err := config.HTTPConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading http config: %w", err)
	}

	jwtCfg, err := config.JWTConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading jwt config: %w", err)
	}

	dbPool, err := pkg.InitDBPool(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("error initializing DB pool: %w", err)
	}

	//init repo
	repo := repository.NewRepository(dbPool)

	// init service
	serv := service.NewService(repo, jwtCfg.Jwt)

	//init router
	v := middleware.NewValidator()
	r := handler.NewRouter(serv, jwtCfg.Jwt, v)

	return &App{
			router:  r,
			httpCfg: htppCfg,
		},
		nil
}

func (a *App) Run() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.httpCfg.GetPort()),
		Handler:      a.router,
		ReadTimeout:  a.httpCfg.GetTimeout(),
		WriteTimeout: a.httpCfg.GetTimeout(),
		IdleTimeout:  a.httpCfg.GetIdleTimeout(),
	}

	// Запуск сервера
	go func() {
		log.Info("Starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server ListenAndServe failed", log.Any("err", err))
		}
	}()

	// Слушаем сигналы остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutdown signal received")

	// Контекст с таймаутом на graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", log.Any("err", err))
		return err
	}

	select {
	case <-ctx.Done():
		log.Warn("Shutdown timeout exceeded")
	default:
		log.Info("Server exited gracefully")
	}

	return nil
}
