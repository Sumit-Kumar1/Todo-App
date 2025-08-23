package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"todoapp/internal/errors"

	// postgres import for driver
	_ "github.com/lib/pq"
)

func newDB(logger *slog.Logger) (*sql.DB, error) {
	ctx := context.Background()
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "")
	dbName := getEnvOrDefault("DB_NAME", "todo")

	if strings.TrimSpace(password) == "" {
		logger.LogAttrs(ctx, slog.LevelError, "DATABASE_URL environment variable not set")
		return nil, errors.ConstError("empty db password")
	}

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlConn)
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
