package cmd

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

func newDB(c context.Context) (*sql.DB, error) {
	var (
		errConn error
		ctx     = context.Background()
	)

	dbOnce.Do(func() {
		slog.InfoContext(ctx, "creating new postgres connection")

		db, err := connectDB(c)
		if err != nil {
			errConn = fmt.Errorf("connectDB: %w", err)
			return
		}

		dbInstance = db
	})

	return dbInstance, errConn
}

func connectDB(c context.Context) (*sql.DB, error) {
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "")
	dbName := getEnvOrDefault("DB_NAME", "todo")

	if strings.TrimSpace(password) == "" {
		return nil, errors.NewConstError("empty db password")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.ErrorContext(c, "sql.Open failed", slog.String("dsn", dsn), slog.String("error", err.Error()))
		return nil, err
	}

	// configure a sane pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.ErrorContext(ctx, "db.Ping failed", slog.String("error", err.Error()))
		return nil, err
	}

	slog.InfoContext(ctx, "postgres connection established",
		slog.String("host", host), slog.String("port", port), slog.String("db", dbName))

	return db, nil
}
