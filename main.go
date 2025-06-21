package main

import (
	"context"
	"log/slog"
	"os"

	"todoapp/cmd"
)

func main() {
	ctx := context.Background()
	if err := cmd.Run(ctx, os.Stdout, nil); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, err.Error())
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "server is stopped!!")
}
