package migrations

import (
	"database/sql"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrate(db *sql.DB, method string) error {
	drv, err := sqlite3.WithInstance(db, &sqlite3.Config{
		MigrationsTable: "todo_migrations",
		DatabaseName:    "todo",
	})
	if err != nil {
		return err
	}

	mig, err := migrate.NewWithDatabaseInstance("file://./internal/migrations", "todo", drv)
	if err != nil {
		return err
	}

	if strings.EqualFold("UP", method) {
		if err := mig.Up(); err.Error() == "no change" {
			return nil
		}

		return err
	}

	return mig.Down()
}
