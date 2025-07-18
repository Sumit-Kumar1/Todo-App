package migrations

import "github.com/sqlitecloud/sqlitecloud-go"

const (
	tasksDown = "DROP TABLE IF EXISTS tasks;"
	tasksUp   = `CREATE TABLE IF NOT EXISTS tasks(
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
		description TEXT,
    done_status BOOLEAN NOT NULL CHECK (done_status IN (0, 1)),
    due_date DATE,
    added_at DATETIME NOT NULL,
    modified_at DATETIME);`
)

type M20241013015650 string

// nolint:revive // unused but need this as method
func (m M20241013015650) up(db *sqlitecloud.SQCloud) error {
	return db.Execute(tasksUp)
}

// nolint:revive // unused but need this as method
func (m M20241013015650) down(db *sqlitecloud.SQCloud) error {
	return db.Execute(tasksDown)
}
