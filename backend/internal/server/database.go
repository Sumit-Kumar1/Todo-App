package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"todoapp/internal/errors"

	// postgres import for driver
	_ "github.com/lib/pq"
)

var (
	dbOnce     sync.Once
	dbInstance *sql.DB
)

func newDB(logger *slog.Logger) (*sql.DB, error) {
	var (
		errConn error
		ctx     = context.Background()
	)

	if dbInstance != nil {
		logger.LogAttrs(ctx, slog.LevelInfo, "using existing db instance")

		return dbInstance, errConn
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "creating new db instance")

	dbOnce.Do(func() {
		host := GetEnvOrDefault("DB_HOST", "localhost")
		port := GetEnvOrDefault("DB_PORT", "5432")
		user := GetEnvOrDefault("DB_USER", "postgres")
		password := GetEnvOrDefault("DB_PASSWORD", "")
		dbName := GetEnvOrDefault("DB_NAME", "todo")

		if strings.TrimSpace(password) == "" {
			logger.LogAttrs(ctx, slog.LevelError, "DATABASE_URL environment variable not set")
			errConn = errors.NewConstError("empty db password")
			return
		}

		psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbName)

		db, err := sql.Open("postgres", psqlConn)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "database error opening connection",
				slog.String("error", err.Error()))
			errConn = err
			return
		}

		if err := db.PingContext(ctx); err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "database error pinging connection")
			errConn = err
			return
		}

		dbInstance = db
	})

	return dbInstance, errConn
}
