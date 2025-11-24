package main

import (
	"context"
	"log/slog"
	"os"

	"todoapp/cmd"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("configs/.env"); err != nil {
		slog.LogAttrs(ctx, slog.LevelInfo, "error while loading envs", slog.String("err", err.Error()))
		slog.LogAttrs(ctx, slog.LevelInfo, "loading system/container env")
	}

	if err := cmd.Run(ctx, os.Stdout, nil); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, err.Error())
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "server is stopped!!")
}
