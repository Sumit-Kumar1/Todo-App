package models

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

// GetLoggerFromCtx will always return a logger, if no logger found in ctx will returns default logger
func GetLoggerFromCtx(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(Logger).(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return logger
}

func GetCorrelationID(ctx context.Context) string {
	corrID, ok := ctx.Value(CorrelationID).(string)
	if !ok {
		return uuid.NewString()
	}

	if strings.TrimSpace(corrID) == "" {
		return uuid.NewString()
	}

	return corrID
}
