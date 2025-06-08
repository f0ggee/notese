package cmd

import (
	"Project2/iteranal"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

var Pool *pgxpool.Pool

func Connect() *pgxpool.Pool {
	var err error

	ctx, cancel := iteranal.Contexte()
	defer cancel()

	_, err = pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Info("pgxpool.ParseConfig err:", err)
		return nil
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL environment variable not set")
		return nil
	}

	Pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("Error connecting to database", err)
		return nil
	}

	if err := Pool.Ping(ctx); err != nil {
		slog.Error("Error pinging database", err)
		return nil

	}

	return Pool

}
