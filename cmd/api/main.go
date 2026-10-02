package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourname/gin-template/internal/health"
	"github.com/yourname/gin-template/internal/platform/config"
	"github.com/yourname/gin-template/internal/platform/database"
	"github.com/yourname/gin-template/internal/platform/logger"
	"github.com/yourname/gin-template/internal/platform/server"
	"github.com/yourname/gin-template/internal/user"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

// run wires dependencies together. No business logic lives here.
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := logger.New(cfg.Env, cfg.LogLevel)
	slog.SetDefault(log)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// --- Features (add new ones here) ---
	userHandler := user.NewHandler(user.NewService(user.NewRepository(db)))
	healthHandler := health.NewHandler(sqlDB)

	router := server.NewRouter(cfg, log, healthHandler, userHandler)
	srv := server.New(cfg, log, router)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return srv.Run(ctx)
}
