package migrations

import "github.com/sqlitecloud/sqlitecloud-go"

const (
	tasksDown = "DROP TABLE IF EXISTS tasks;"
	tasksUp   = `CREATE TABLE IF NOT EXISTS tasks(
    task_id TEXT PRIMARY KEY, 
    user_id TEXT NOT NULL,
    task_title TEXT NOT NULL, 
    done_status BOOLEAN NOT NULL CHECK (done_status IN (0, 1)),
    added_at DATETIME NOT NULL, 
    modified_at DATETIME);`
)

type M20241013015650 string

func (m M20241013015650) Up(db *sqlitecloud.SQCloud) error {
	return run(db, tasksUp, m)
}

func (m M20241013015650) Down(db *sqlitecloud.SQCloud) error {
	return run(db, tasksDown, m)
}
