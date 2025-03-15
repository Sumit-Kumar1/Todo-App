package models

import (
	"context"
	"log/slog"
)

type CtxLogger string

const Logger CtxLogger = "logger"

// GetLoggerFromCtx will always return a logger, if no logger found in ctx will returns default logger
func GetLoggerFromCtx(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(Logger).(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return logger
}
