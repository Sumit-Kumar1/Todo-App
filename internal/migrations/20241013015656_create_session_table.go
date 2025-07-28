package migrations

import "database/sql"

const (
	sessionDown = "DROP TABLE IF EXISTS sessions;"
	sessionUp   = `CREATE TABLE IF NOT EXISTS sessions(
    id VARCHAR(50) PRIMARY KEY, 
    user_id VARCHAR(50) NOT NULL UNIQUE,
    token TEXT NOT NULL UNIQUE, 
    expiry TIMESTAMP NOT NULL);`
)

type M20241013015656 string

// nolint:revive // unused but need this as method
func (m M20241013015656) up(db *sql.Tx) error {
	_, err := db.Exec(sessionUp)
	return err
}

// nolint:revive // unused but need this as method
func (m M20241013015656) down(db *sql.Tx) error {
	_, err := db.Exec(sessionDown)
	return err
}
