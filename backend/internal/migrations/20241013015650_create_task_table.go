package migrations

import "database/sql"

const (
	tasksDown = "DROP TABLE IF EXISTS tasks;"
	taskTable = `CREATE TABLE IF NOT EXISTS tasks(
    id VARCHAR(50) PRIMARY KEY,
	parent_id VARCHAR(36),
    user_id VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
	description TEXT,
    status VARCHAR(50),
    due_date TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
	FOREIGN KEY (parent_id) REFERENCES tasks(id) ON DELETE CASCADE);`
)

type M20241013015650 string

// nolint:revive // unused but need this as method
func (m M20241013015650) up(db *sql.Tx) error {
	_, err := db.Exec(taskTable)
	return err
}

// nolint:revive // unused but need this as method
func (m M20241013015650) down(db *sql.Tx) error {
	_, err := db.Exec(tasksDown)
	return err
}
