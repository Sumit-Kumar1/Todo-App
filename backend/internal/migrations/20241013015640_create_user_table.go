package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const (
	userTable = `CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(50) NOT NULL PRIMARY KEY, 
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password BYTEA NOT NULL);`
)

func createTableUser() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(userTable)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
