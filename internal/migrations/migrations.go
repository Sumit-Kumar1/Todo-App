package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"todoapp/internal/models"
	"todoapp/internal/server"
)

const (
	methodUp       = "UP"
	methodDown     = "DOWN"
	migTableName   = "todo_migrations"
	migInsertErr   = "Migration table insert error"
	createMigTable = "CREATE TABLE IF NOT EXISTS %s(version TEXT, start_time DATETIME, end_time DATETIME, method TEXT);"
	versionQuery   = "SELECT version from ? ORDER BY version DESC"
	insertVersion  = "INSERT INTO %s(version, start_time, method) VALUES (?, ?,?);"
)

type migrator interface {
	up(tx *sql.Tx) error
	down(tx *sql.Tx) error
}

func RunMigrations(ctx context.Context, s *server.Server, method string) error {
	if s.DB == nil {
		return models.NewConstError("db is nil")
	}

	t := time.Now()

	query := fmt.Sprintf(createMigTable, migTableName)

	_, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "not able to create the migration table")

		return err
	}

	switch method {
	case methodUp:
		err = runUpMigrations(ctx, s, migrations)
	case methodDown:
		err = runDownMigrations(ctx, s, migrations)
	default:
		return models.ErrInvalid("migration method")
	}

	if err != nil {
		return err
	}

	s.Logger.LogAttrs(ctx, slog.LevelInfo,
		fmt.Sprintf("completed the migration in time: %v seconds", time.Since(t).Seconds()),
	)

	return nil
}

func runUpMigrations(ctx context.Context, s *server.Server, migs map[string]migrator) error {
	run := make([]string, 0)

	lastRun, err := getLastRunMigration(ctx, s)
	if err != nil {
		return err
	}

	for version, val := range migs {
		if version <= lastRun {
			continue
		}

		if err := performUpMigrations(ctx, s, val, version); err != nil {
			return err
		}

		run = append(run, version)
	}

	s.Logger.LogAttrs(ctx, slog.LevelInfo, "successfully UP migrated",
		slog.String("runs", fmt.Sprintf("[%s]", strings.Join(run, ", "))))

	return nil
}

func runDownMigrations(ctx context.Context, s *server.Server, migs map[string]migrator) error {
	run := []string{}
	versions := []string{}
	version := ""

	rows, err := s.DB.QueryContext(ctx, versionQuery, migTableName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.LogAttrs(ctx, slog.LevelWarn, "no versions found to revert")

			return nil
		}

		return err
	}

	for rows.Next() {
		if err := rows.Scan(&version); err != nil {
			return err
		}

		versions = append(versions, version)
	}

	for idx := range versions {
		mig := migs[versions[idx]]

		if err := performDownMigrations(ctx, s, mig, versions[idx]); err != nil {
			return err
		}

		run = append(run, versions[idx])
	}

	s.Logger.LogAttrs(ctx, slog.LevelInfo, "successfully 'DOWN' migrated",
		slog.String("runs", fmt.Sprintf("[%s]", strings.Join(run, ", "))))

	return nil
}

func getLastRunMigration(ctx context.Context, s *server.Server) (string, error) {
	var (
		lastRun         string
		queryGetLastRun = fmt.Sprintf("SELECT version FROM %s ORDER BY version DESC LIMIT 1;", migTableName)
	)

	res, err := s.DB.QueryContext(ctx, queryGetLastRun)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogAttrs(ctx, slog.LevelWarn, "no last run migrations found")

			return "", nil
		}

		return "", err
	}

	for res.Next() {
		if err := res.Scan(&lastRun); err != nil {
			s.Logger.LogAttrs(ctx, slog.LevelError, "unable to fetch last run of the migration",
				slog.String("error", err.Error()),
			)

			return "", err
		}
	}

	return lastRun, nil
}

func performUpMigrations(ctx context.Context, s *server.Server, val migrator, version string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "unable to start transaction",
			slog.String("error", err.Error()),
		)

		return err
	}

	query := fmt.Sprintf(insertVersion, migTableName)

	_, err = tx.ExecContext(ctx, query, version, time.Now().UnixMilli(), methodUp)
	if err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, migInsertErr,
			slog.String("migration", version),
			slog.String("error", err.Error()),
		)

		return handleRollback(tx, err)
	}

	if err := val.up(tx); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "Migration error",
			slog.String("migration", version),
			slog.String("error", err.Error()),
		)

		return handleRollback(tx, err)
	}

	query = fmt.Sprintf("UPDATE %s SET end_time=? WHERE version=?;", migTableName)

	if _, err := tx.ExecContext(ctx, query, time.Now().UnixMilli(), version); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, migInsertErr,
			slog.String("migration", version),
			slog.String("error", err.Error()),
		)

		return handleRollback(tx, err)
	}

	return tx.Commit()
}

func performDownMigrations(ctx context.Context, s *server.Server, val migrator, key string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "unable to start transaction",
			slog.String("error", err.Error()),
		)

		return err
	}

	if err := val.down(tx); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "Migration error",
			slog.String("migration", key),
			slog.String("error", err.Error()),
		)

		return handleRollback(tx, err)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE version=?", migTableName)

	if _, err := tx.ExecContext(ctx, query, key); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, migInsertErr,
			slog.String("migration", key),
			slog.String("error", err.Error()),
		)

		return handleRollback(tx, err)
	}

	return tx.Commit()
}

func handleRollback(tx *sql.Tx, err error) error {
	if rErr := tx.Rollback(); rErr != nil {
		return rErr
	}

	return err
}
