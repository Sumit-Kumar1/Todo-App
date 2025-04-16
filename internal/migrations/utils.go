package migrations

import (
	"fmt"
	"log/slog"

	"github.com/sqlitecloud/sqlitecloud-go"
)

func run(db *sqlitecloud.SQCloud, query string, m migrator) error {
	err := db.Execute(query)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("run migration: %s", m))

	return nil
}
