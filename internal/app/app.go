package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shuyou-ai/shuyou-go/internal/config"
	"github.com/shuyou-ai/shuyou-go/internal/handler"
	"github.com/shuyou-ai/shuyou-go/internal/infra/database"
	"github.com/shuyou-ai/shuyou-go/internal/infra/jwt"
	"github.com/shuyou-ai/shuyou-go/internal/infra/logger"
	"github.com/shuyou-ai/shuyou-go/internal/repository"
	"github.com/shuyou-ai/shuyou-go/internal/router"
	"github.com/shuyou-ai/shuyou-go/internal/service"
	pkgvalidator "github.com/shuyou-ai/shuyou-go/pkg/validator"
	"go.uber.org/zap"
)

type App struct {
	cfg    *config.Config
	log    *zap.Logger
	db     *database.Client
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	log, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	pkgvalidator.Init()

	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.EnsureIndexes(ctx); err != nil {
		_ = db.Close(context.Background())
		return nil, fmt.Errorf("ensure indexes: %w", err)
	}

	userRepo := repository.NewUserRepository(db.DB())
	userService := service.NewUserService(userRepo)
	jwtManager := jwt.NewManager(cfg.JWT)

	engine := router.New(cfg.Server.Mode, router.Deps{
		Log: log,
		JWT: jwtManager,
		Handlers: &router.Handlers{
			Health: handler.NewHealthHandler(db),
			User:   handler.NewUserHandler(userService, jwtManager),
		},
	})

	server := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &App{cfg: cfg, log: log, db: db, server: server}, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 1)

	go func() {
		a.log.Info("server starting",
			zap.String("addr", a.server.Addr),
			zap.String("mode", a.cfg.Server.Mode),
		)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-quit:
		a.log.Info("server shutting down", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	if a.db != nil {
		if err := a.db.Close(ctx); err != nil {
			return fmt.Errorf("close mongodb: %w", err)
		}
	}

	_ = a.log.Sync()
	return nil
}
