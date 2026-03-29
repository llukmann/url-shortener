package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	config := config.MustLoad()

	log := setupLogger(config.Env)

	log.Info("starting server", slog.String("env", config.Env))
	log.Debug("Debug level is enabled")

	ctxInitDb, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	storage, err := postgres.NewStorage(ctxInitDb, createDatabaseUrl(&config.Database))
	if err != nil {
		log.Error("failed to init storage", err.Error())
		os.Exit(1)
	}

	ctxCloseDb := context.Background()
	err = storage.Close(ctxCloseDb)
	if err != nil {
		log.Error("failed to close storage", err.Error())
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}

func createDatabaseUrl(dbConfig *config.Database) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Dbname,
		dbConfig.Sslmode,
	)
}
