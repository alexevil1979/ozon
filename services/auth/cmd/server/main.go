// Package main — Auth service entrypoint. Echo REST + health + PostgreSQL (pgx) + JWT.
// Choice: Echo over Fiber for stdlib-aligned middleware and OWASP-friendly defaults.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"marketplace/auth/internal/auth"
	"marketplace/auth/internal/config"
	"marketplace/auth/internal/db"
	"marketplace/auth/internal/handler"
	"marketplace/auth/internal/jwt"
	"marketplace/auth/internal/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("db connect", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		slog.Error("migrate", "err", err)
		os.Exit(1)
	}

	jwtUtil := jwt.New(cfg.JWTSecret, cfg.JWTExpiry)
	refreshExpiry, _ := time.ParseDuration(cfg.RefreshExpiry)
	if refreshExpiry == 0 {
		refreshExpiry = 7 * 24 * time.Hour
	}
	userRepo := repository.NewUserRepository(pool)
	refreshTokenRepo := repository.NewRefreshTokenRepository(pool)
	authSvc := auth.NewService(userRepo, refreshTokenRepo, jwtUtil, refreshExpiry)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))

	handler.Register(e, pool, authSvc)

	addr := ":" + cfg.Port
	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			slog.Error("server", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown", "err", err)
	}
	slog.Info("auth service stopped")
}
