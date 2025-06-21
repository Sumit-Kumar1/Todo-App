package migrations

import "github.com/sqlitecloud/sqlitecloud-go"

const (
	sessionDown = "DROP TABLE IF EXISTS sessions;"
	sessionUp   = `CREATE TABLE IF NOT EXISTS sessions(
    id TEXT PRIMARY KEY, 
    user_id TEXT NOT NULL UNIQUE,
    token TEXT NOT NULL UNIQUE, 
    expiry DATETIME NOT NULL);`
)

type M20241013015656 string

// nolint:revive // unused but need this as method
func (m M20241013015656) up(db *sqlitecloud.SQCloud) error {
	return db.Execute(sessionUp)
}

// nolint:revive // unused but need this as method
func (m M20241013015656) down(db *sqlitecloud.SQCloud) error {
	return db.Execute(sessionDown)
}
