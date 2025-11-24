package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"todoapp/internal/errors"

	_ "github.com/lib/pq"
)

var (
	dbOnce     sync.Once
	dbInstance *sql.DB
)

func newDB() (*sql.DB, error) {
	var (
		errConn error
		ctx     = context.Background()
	)

<<<<<<< HEAD
	if dbInstance != nil {
		logger.LogAttrs(ctx, slog.LevelDebug, "using existing db instance")

		return dbInstance, errConn
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "creating new db instance")

=======
>>>>>>> f7dd975392b109c2deb927d548e51c8af7f4cad1
	dbOnce.Do(func() {
		slog.InfoContext(ctx, "creating new postgres connection")

<<<<<<< HEAD
		if strings.TrimSpace(password) == "" {
			logger.LogAttrs(ctx, slog.LevelError, "DB_PASSWORD environment variable not set")
			errConn = errors.NewConstError("empty db password")
			return
		}

		psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbName)

		db, err := sql.Open("postgres", psqlConn)
=======
		db, err := connectDB(ctx)
>>>>>>> f7dd975392b109c2deb927d548e51c8af7f4cad1
		if err != nil {
			errConn = fmt.Errorf("connectDB: %w", err)
			return
		}

		dbInstance = db
	})

	return dbInstance, errConn
}

func connectDB(ctx context.Context) (*sql.DB, error) {
	host := GetEnvOrDefault("DB_HOST", "localhost")
	port := GetEnvOrDefault("DB_PORT", "5432")
	user := GetEnvOrDefault("DB_USER", "postgres")
	password := GetEnvOrDefault("DB_PASSWORD", "")
	dbName := GetEnvOrDefault("DB_NAME", "todo")

	if strings.TrimSpace(password) == "" {
		return nil, errors.NewConstError("empty db password")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.ErrorContext(ctx, "sql.Open failed", slog.String("dsn", dsn), slog.String("error", err.Error()))
		return nil, err
	}

	// configure a sane pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.ErrorContext(ctx, "db.Ping failed", slog.String("error", err.Error()))
		return nil, err
	}

	slog.InfoContext(ctx, "postgres connection established",
		slog.String("host", host), slog.String("port", port), slog.String("db", dbName))

	return db, nil
}
