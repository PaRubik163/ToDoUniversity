package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"homework-reminder-bot/internal/app"
	"homework-reminder-bot/internal/config"
	"homework-reminder-bot/internal/database"
	"homework-reminder-bot/internal/ratelimit"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	limiter, err := ratelimit.NewRedis(cfg.RedisURL)
	if err != nil {
		return err
	}
	defer limiter.Close()

	return app.New(cfg, db, limiter).Run(ctx)
}
