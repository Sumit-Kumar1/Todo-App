package server

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/sqlitecloud/sqlitecloud-go"
)

func newDB(logger *slog.Logger) (*sqlitecloud.SQCloud, error) {
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
		logger.Error(err.Error())
		return nil, err
	}

	if !sqcl.IsConnected() {
		return nil, fmt.Errorf("not able to connect to database")
	}

	logger.Info("DB connected successfully")

	return sqcl, nil
}
