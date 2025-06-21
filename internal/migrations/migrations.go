package migrations

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"todoapp/internal/models"
	"todoapp/internal/server"

	"github.com/sqlitecloud/sqlitecloud-go"
)

const (
	migTableName = "todo_migrations"
	migInsertErr = "Migration table insert error"
)

type migrator interface {
	up(db *sqlitecloud.SQCloud) error
	down(db *sqlitecloud.SQCloud) error
}

func RunMigrations(ctx context.Context, s *server.Server, method string) error {
	if s.DB == nil {
		return models.NewConstError("db is nil")
	}

	t := time.Now()

	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(version TEXT, start_time DATETIME, end_time DATETIME, method TEXT);",
		migTableName,
	)

	err := s.DB.Execute(query)
	if err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "not able to create the migration table")

		return err
	}

	switch method {
	case "UP":
		err = runUpMigrations(ctx, s, migrations)
	case "DOWN":
		err = runDownMigrations(ctx, s, migrations)
	default:
		return models.ErrInvalid("migration method")
	}

	if err != nil {
		return err
	}

	s.Logger.LogAttrs(ctx, slog.LevelInfo,
		fmt.Sprintf("Completed the migration in time: %v seconds", time.Since(t).Seconds()),
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

	getAllVersions := fmt.Sprintf("SELECT version from %s ORDER BY version DESC", migTableName)

	rows, err := s.DB.Select(getAllVersions)
	if err != nil {
		return err
	}

	numRows := rows.GetNumberOfRows()
	for r := uint64(0); r < numRows; r++ {
		var version string

		version, err := rows.GetStringValue(r, 0)
		if err != nil {
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

	res, err := s.DB.Select(queryGetLastRun)
	if err != nil {
		return "", err
	}

	numRows := res.GetNumberOfRows()
	if numRows == 0 {
		s.Logger.LogAttrs(ctx, slog.LevelInfo, "no last run migrations found")

		return "", nil
	}

	for i := uint64(0); i < res.GetNumberOfRows(); i++ {
		lastRun, err = res.GetStringValue(i, 0)
		if err != nil {
			s.Logger.LogAttrs(ctx, slog.LevelError, "unable to fetch last run of the migration",
				slog.String("error", err.Error()),
			)

			return "", err
		}
	}

	return lastRun, nil
}

func performUpMigrations(ctx context.Context, s *server.Server, val migrator, version string) error {
	const method = "UP"

	if err := s.DB.BeginTransaction(); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "unable to start transaction",
			slog.String("error", err.Error()),
		)

		return err
	}

	defer func() {
		_ = s.DB.EndTransaction()
	}()

	query := fmt.Sprintf("INSERT INTO %s (version, start_time, method) VALUES ('%s', %v,'%s');",
		migTableName, version, time.Now().UnixMilli(), method)

	if err := s.DB.Execute(query); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, migInsertErr,
			slog.String("migration", version),
			slog.String("error", err.Error()),
		)

		return handleRollback(s, err)
	}

	if err := val.up(s.DB); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "Migration error",
			slog.String("migration", version),
			slog.String("error", err.Error()),
		)

		return handleRollback(s, err)
	}

	query = fmt.Sprintf("UPDATE %s SET end_time=%v WHERE version='%s';",
		migTableName,
		time.Now().UnixMilli(),
		version,
	)

	if err := s.DB.Execute(query); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, migInsertErr,
			slog.String("migration", version),
			slog.String("error", err.Error()),
		)

		return handleRollback(s, err)
	}

	return nil
}

func performDownMigrations(ctx context.Context, s *server.Server, val migrator, key string) error {
	if err := s.DB.BeginTransaction(); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "unable to start transaction",
			slog.String("error", err.Error()),
		)

		return err
	}

	defer func() {
		_ = s.DB.EndTransaction()
	}()

	if err := val.down(s.DB); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, "Migration error",
			slog.String("migration", key),
			slog.String("error", err.Error()),
		)

		return handleRollback(s, err)
	}

	if err := s.DB.Execute(fmt.Sprintf("DELETE FROM %s WHERE version = %v", migTableName, key)); err != nil {
		s.Logger.LogAttrs(ctx, slog.LevelError, migInsertErr,
			slog.String("migration", key),
			slog.String("error", err.Error()),
		)

		return handleRollback(s, err)
	}

	return nil
}

func handleRollback(s *server.Server, err error) error {
	if rErr := s.DB.RollBackTransaction(); rErr != nil {
		return rErr
	}

	return err
}
