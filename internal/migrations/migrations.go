package migrations

import (
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
	migTableName = "todo_migrations"
)

type Migrator interface {
	Up(db *sql.Tx) error
	Down(db *sql.Tx) error
}

func RunMigrations(s *server.Server, method string) error {
	if s.DB == nil {
		return fmt.Errorf("db is nil")
	}

	t := time.Now()

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(version TEXT, start_time DATETIME, end_time DATETIME, method TEXT);",
		migTableName)

	_, err := s.DB.Exec(query)
	if err != nil {
		s.Logger.Error("not able to create the migration table", slog.String("error", err.Error()))
		return err
	}

	switch method {
	case "UP":
		err = runUpMigrations(s, migrations)
	case "DOWN":
		err = runDownMigrations(s, migrations)
	default:
		s.Logger.Error("invalid migrations method provided!!")
		return models.ErrInvalid("migration method")
	}

	if err != nil {
		return err
	}

	s.Logger.Info(fmt.Sprintf("Completed the migration in time: %v seconds", time.Since(t).Seconds()))

	return nil
}

func runUpMigrations(s *server.Server, migs map[string]Migrator) error {
	var run []string

	lastRun, err := getLastRunMigration(s)
	if err != nil {
		return err
	}

	for key, val := range migs {
		if key <= lastRun {
			continue
		}

		if err := performMigration(s, val, key, "UP"); err != nil {
			return err
		}

		run = append(run, key)
	}

	s.Logger.Info("successfully UP migrated",
		slog.String("runs", fmt.Sprintf("[%s]", strings.Join(run, ", "))))

	return nil
}

func runDownMigrations(s *server.Server, migs map[string]Migrator) error {
	var (
		run      []string
		versions []string
	)

	getAllVersions := fmt.Sprintf("SELECT version from %s ORDER BY version DESC", migTableName)
	rows, err := s.DB.Query(getAllVersions)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			s.Logger.Error("failed to scan row", slog.String("error", err.Error()))
			return err
		}

		versions = append(versions, version)
	}

	for idx := range versions {
		migr := migs[versions[idx]]

		if err := performMigration(s, migr, versions[idx], "DOWN"); err != nil {
			return err
		}

		run = append(run, versions[idx])
	}

	s.Logger.Info("successfully 'DOWN' migrated",
		slog.String("runs", fmt.Sprintf("[%s]", strings.Join(run, ", "))))

	return nil
}

func getLastRunMigration(s *server.Server) (string, error) {
	var (
		lastRun         string
		queryGetLastRun = fmt.Sprintf("SELECT version FROM %s ORDER BY version DESC LIMIT 1;", migTableName)
	)

	err := s.DB.QueryRow(queryGetLastRun).Scan(&lastRun)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.Info("no last run migrations found")
			return "", nil
		}

		s.Logger.Error("not able to get the last run of the migration", slog.String("error", err.Error()))
		return "", err
	}

	return lastRun, nil
}

func performMigration(s *server.Server, val Migrator, key, method string) error {
	var (
		upPreRun    = fmt.Sprintf("INSERT INTO %s(version, start_time, method) VALUES (?, ?, ?)", migTableName)
		upPostRun   = fmt.Sprintf("UPDATE %s SET end_time=? WHERE version=?", migTableName)
		downPostRun = fmt.Sprintf("DELETE FROM %s WHERE version = ?", migTableName)
	)

	tx, err := s.DB.Begin()
	if err != nil {
		s.Logger.Error("unable to start transaction", slog.String("error", err.Error()))
		return err
	}

	switch method {
	case "UP":
		if _, err := tx.Exec(upPreRun, key, time.Now(), method); err != nil {
			s.Logger.Error("Migration table insert error", slog.String("migration", key), slog.String("error", err.Error()))

			return handleRollback(tx, err)
		}

		if err := val.Up(tx); err != nil {
			s.Logger.Error("Migration error", slog.String("migration", key), slog.String("error", err.Error()))

			return handleRollback(tx, err)
		}

		if _, err := tx.Exec(upPostRun, time.Now(), key); err != nil {
			s.Logger.Error("Migration table insert error", slog.String("migration", key), slog.String("error", err.Error()))

			return handleRollback(tx, err)
		}

	case "DOWN":
		if err := val.Down(tx); err != nil {
			s.Logger.Error("Migration error", slog.String("migration", key), slog.String("error", err.Error()))

			return handleRollback(tx, err)
		}

		if _, err := tx.Exec(downPostRun, key); err != nil {
			s.Logger.Error("Migration table insert error", slog.String("migration", key), slog.String("error", err.Error()))

			return handleRollback(tx, err)
		}

	default:
		s.Logger.Error("invalid migration method provided!!")
		return models.ErrInvalid("migration method")
	}

	if err := tx.Commit(); err != nil {
		s.Logger.Error("unable to commit transaction", slog.String("migration", key), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func handleRollback(tx *sql.Tx, err error) error {
	if rErr := tx.Rollback(); rErr != nil {
		return rErr
	}

	return err
}
