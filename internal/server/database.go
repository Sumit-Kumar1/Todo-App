package server

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"strings"
	"todoapp/internal/models"

	// Used this for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func newDB(logger *slog.Logger) (*sql.DB, error) {
	ctx := context.Background()

	var dbURLErr = models.NewConstError("DATABASE_URL environment variable not set")

	database := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(database) == "" {
		logger.LogAttrs(ctx, slog.LevelError, "DATABASE_URL environment variable not set")
		return nil, dbURLErr
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
