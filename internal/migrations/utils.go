package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

func run(db *sql.Tx, query string, m Migrator) error {
	_, err := db.ExecContext(context.Background(), query)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("run migration: %s", m))

	return nil
}
