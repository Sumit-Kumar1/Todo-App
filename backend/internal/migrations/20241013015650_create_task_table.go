package migrations

import "gofr.dev/pkg/gofr/migration"

const (
	taskTable = `CREATE TABLE IF NOT EXISTS tasks(
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
	description TEXT,
    done_status BOOLEAN,
    due_date TIMESTAMP,
    added_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP);`
)

func createTableTask() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(taskTable)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
