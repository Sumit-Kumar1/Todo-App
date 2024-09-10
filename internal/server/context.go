package server

import (
	"context"
	"database/sql"
	"html/template"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

type DataStore struct {
	DB *sql.DB
	*DBConfig
}

type DBConfig struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         string
	Dialect      string
}

type Context struct {
	*DataStore
	Template *template.Template
	Logger   *slog.Logger
	context.Context
}

func NewContext() *Context {
	logger := newLogger()
	templ := newTemplate()

	db, err := newDB(logger)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return &Context{
		Logger:   logger,
		DB:       db,
		Template: templ,
	}
}

func (c *Context) Cleanup() error {
	if err := c.DB.Close(); err != nil {
		c.Logger.Error(err.Error())
		return err
	}

	ch := c.Context.Done()
	if ch == nil {
		c.Logger.Info("context can't be done yet")
		return c.Err()
	}

	return nil
}

func newDB(logger *slog.Logger) (*sql.DB, error) {
	const dbFile string = "../tasks.db"

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		logger.Error("unable to connect sqlite", "error", err.Error())
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Error("database not reachable", "error", err)
		return nil, err
	}

	return db, nil
}

func newLogger() *slog.Logger {
	var (
		level slog.Level
	)

	switch os.Getenv("LOG_LEVEL") {
	case "WARN":
		level = slog.LevelWarn
	case "DEBUG":
		level = slog.LevelDebug
	case "Error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))

	slog.SetDefault(logger)

	return logger
}

func newTemplate() *template.Template {
	pattern := "views/*"

	return template.Must(template.ParseGlob(pattern))
}
