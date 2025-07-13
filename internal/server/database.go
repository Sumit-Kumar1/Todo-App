package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func newDB(logger *slog.Logger) (*sql.DB, error) {
	ctx := context.Background()

	database := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(database) == "" {
		logger.LogAttrs(ctx, slog.LevelError, "DATABASE_URL environment variable not set")
		return nil, errors.New("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "database error opening connection",
			slog.String("error", err.Error()))
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "database error pinging connection")
		return nil, err
	}

	return db, nil
}
