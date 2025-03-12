package server

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"todoapp/internal/models"

	"github.com/sqlitecloud/sqlitecloud-go"
)

func newDB(logger *slog.Logger) (*sqlitecloud.SQCloud, error) {
	ctx := context.Background()

	config := sqlitecloud.SQCloudConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     getEnvAsInt("DB_PORT", 8860),
		Database: os.Getenv("DB_NAME"),
		ApiKey:   os.Getenv("DB_APIKEY"),
		MaxRows:  getEnvAsInt("DB_MAXROWS", 20),
	}

	isSecure, err := strconv.ParseBool(os.Getenv("DB_SECURE_FLAG"))
	if err != nil {
		return nil, err
	}

	config.Secure = isSecure

	sqcl := sqlitecloud.New(config)

	if err := sqcl.Connect(); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error while connecting to Database", slog.String("error", err.Error()))
		return nil, err
	}

	if !sqcl.IsConnected() {
		return nil, models.NewConstError("database is not connected after conn success")
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "DB connected successfully")

	return sqcl, nil
}
