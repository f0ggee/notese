package cmd

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

func Connect(ctx context.Context) (*sql.DB, error) {

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("Error loading database URL")
		return nil, errors.New("Error loading database URL")

	}

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		slog.Error("Error loading database connection")
		os.Exit(1)
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		slog.Error("Error pinging database")
		os.Exit(1)
		return nil, err
	}

	slog.Info("Successfully connected to database")
	return db, nil

}
