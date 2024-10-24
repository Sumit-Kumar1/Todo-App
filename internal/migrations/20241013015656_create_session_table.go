package migrations

import "database/sql"

const (
	sessionDown = "DROP TABLE IF EXISTS sessions;"
	sessionUp   = `CREATE TABLE IF NOT EXISTS sessions(
    id TEXT PRIMARY KEY, 
    user_id TEXT NOT NULL UNIQUE,
    token TEXT NOT NULL UNIQUE, 
    expiry DATETIME NOT NULL);`
)

type M20241013015656 string

func (m M20241013015656) Up(db *sql.Tx) error {
	return run(db, sessionUp, m)
}

func (m M20241013015656) Down(db *sql.Tx) error {
	return run(db, sessionDown, m)
}
