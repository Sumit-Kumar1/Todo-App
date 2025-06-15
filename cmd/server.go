package cmd

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"todoapp/internal/migrations"
	"todoapp/internal/models"
	"todoapp/internal/server"
)

func Run(c context.Context, _ io.Writer, _ []string) error {
	ctx, stop := signal.NotifyContext(c, os.Interrupt)
	defer stop()

	app, err := server.NewServer()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// add logger into context
	ctx = context.WithValue(ctx, models.Logger, app.Logger)

	server.SetupRoutes(ctx, app)

	if err = migrations.RunMigrations(ctx, app, app.MigrationMethod); err != nil {
		slog.LogAttrs(c, slog.LevelError, "error while running migrations",
			slog.String("error", err.Error()))

		return err
	}

	srvErr := make(chan error, 1)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(app.Host, app.Port),
		Handler:      app.GlobalRateLimiter(app.Mux),
		ReadTimeout:  time.Duration(app.ReadTimeout * int(time.Second)),
		WriteTimeout: time.Duration(app.WriteTimeout * int(time.Second)),
		IdleTimeout:  time.Duration(app.IdleTimeout * int(time.Second)),
	}

	go func() {
		app.Logger.LogAttrs(ctx, slog.LevelInfo, "Server started", slog.String("Address", httpServer.Addr))

		srvErr <- httpServer.ListenAndServe()
	}()

	if err := checkForTrigger(ctx, app, srvErr); err != nil {
		return err
	}

	if err := httpServer.Shutdown(context.Background()); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "error while shutting down the server",
			slog.String("error", err.Error()),
		)
	}

	return nil
}

func checkForTrigger(ctx context.Context, app *server.Server, srvErr chan error) error {
	var err error

	select {
	case err = <-srvErr:
		if !errors.Is(err, http.ErrServerClosed) {
			app.Logger.LogAttrs(ctx, slog.LevelError,
				"error listening and serving",
				slog.String("error", err.Error()),
			)

			return err
		}

		if app.DB == nil {
			return nil
		}

		if err := app.DB.Close(); err != nil {
			app.Logger.LogAttrs(ctx, slog.LevelError, "error while closing database connection",
				slog.String("error", err.Error()))

			return err
		}

		app.Logger.LogAttrs(ctx, slog.LevelInfo, "database connection closed successfully")

	case <-ctx.Done():
	}

	return nil
}
