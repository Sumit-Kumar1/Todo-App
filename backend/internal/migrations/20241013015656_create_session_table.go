package migrations

import "gofr.dev/pkg/gofr/migration"

const (
	sessionTable = `CREATE TABLE IF NOT EXISTS sessions(
    id VARCHAR(50) PRIMARY KEY, 
    user_id VARCHAR(50) NOT NULL UNIQUE,
    token TEXT NOT NULL UNIQUE, 
    expiry TIMESTAMP NOT NULL);`
)

func createTableSession() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(sessionTable)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
