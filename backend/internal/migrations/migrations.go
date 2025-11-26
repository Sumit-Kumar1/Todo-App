package migrations

import (
	"context"
	"database/sql"
	liberrors "errors"
	"fmt"
	"log/slog"
	"time"

	"todoapp/internal/errors"
)

const (
	methodUp       = "UP"
	methodDown     = "DOWN"
	migTableName   = "todo_migrations"
	migInsertErr   = "Migration table insert error"
	createMigTable = "CREATE TABLE IF NOT EXISTS %s(version TEXT, start_time TIMESTAMP, end_time TIMESTAMP, method TEXT);"
	versionQuery   = "SELECT version from %s ORDER BY version DESC"
	insertVersion  = "INSERT INTO %s(version, start_time, method) VALUES ($1, $2, $3);"
)

type migrator interface {
	up(tx *sql.Tx) error
	down(tx *sql.Tx) error
}

func RunMigrations(ctx context.Context, db *sql.DB, method string) error {
	if db == nil {
		return errors.NewConstError("db is nil")
	}

	query := fmt.Sprintf(createMigTable, migTableName)

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	switch method {
	case methodUp:
		err = runUpMigrations(ctx, db, migrations)
	case methodDown:
		err = runDownMigrations(ctx, db, migrations)
	default:
		return errors.ErrInvalid("migration method")
	}

	if err != nil {
		return err
	}

	return nil
}

func runUpMigrations(ctx context.Context, s *sql.DB, migs map[string]migrator) error {
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

	return nil
}

func runDownMigrations(ctx context.Context, db *sql.DB, migs map[string]migrator) error {
	var (
		run      []string
		versions []string
		version  string
	)

	query := fmt.Sprintf(versionQuery, migTableName)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		if liberrors.Is(err, sql.ErrNoRows) {
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

		if err := performDownMigrations(ctx, db, mig, versions[idx]); err != nil {
			return err
		}

		run = append(run, versions[idx])
	}

	return nil
}

func getLastRunMigration(ctx context.Context, db *sql.DB) (string, error) {
	var (
		lastRun         string
		queryGetLastRun = fmt.Sprintf("SELECT version FROM %s ORDER BY version DESC LIMIT 1;", migTableName)
	)

	res, err := db.QueryContext(ctx, queryGetLastRun)
	if err != nil {
		if liberrors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	for res.Next() {
		if err := res.Scan(&lastRun); err != nil {
			return "", err
		}
	}

	return lastRun, nil
}

func performUpMigrations(ctx context.Context, db *sql.DB, val migrator, version string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(insertVersion, migTableName)

	_, err = tx.ExecContext(ctx, query, version, time.Now(), methodUp)
	if err != nil {
		return handleRollback(tx, err)
	}

	if err := val.up(tx); err != nil {
		return handleRollback(tx, err)
	}

	query = fmt.Sprintf("UPDATE %s SET end_time=$1 WHERE version=$2;", migTableName)

	if _, err := tx.ExecContext(ctx, query, time.Now(), version); err != nil {
		return handleRollback(tx, err)
	}

	return tx.Commit()
}

func performDownMigrations(ctx context.Context, db *sql.DB, val migrator, key string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := val.down(tx); err != nil {
		return handleRollback(tx, err)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE version=$1", migTableName)

	if _, err := tx.ExecContext(ctx, query, key); err != nil {
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
