package migrations

import (
	"database/sql"
)

const (
	tasksDown = "DROP TABLE IF EXISTS tasks;"
	tasksUp   = `CREATE TABLE IF NOT EXISTS tasks(
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
	description TEXT,
    done_status BOOLEAN,
    due_date TIMESTAMP,
    added_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP);`
)

type M20241013015650 string

// nolint:revive // unused but need this as method
func (m M20241013015650) up(db *sql.Tx) error {
	_, err := db.Exec(tasksUp)
	return err
}

// nolint:revive // unused but need this as method
func (m M20241013015650) down(db *sql.Tx) error {
	_, err := db.Exec(tasksDown)
	return err
}
