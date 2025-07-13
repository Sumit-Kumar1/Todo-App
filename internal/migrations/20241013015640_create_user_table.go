package migrations

import (
	"database/sql"
)

const (
	userDown = "DROP TABLE IF EXISTS users;"
	userUp   = `CREATE TABLE IF NOT EXISTS users(
    id TEXT NOT NULL PRIMARY KEY, 
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE CHECK (email LIKE '%'),
    password TEXT NOT NULL);`
)

type M20241013015640 string

// nolint:revive // unused but need this as method
func (m M20241013015640) up(db *sql.Tx) error {
	_, err := db.Exec(userUp)
	return err
}

// nolint:revive // unused but need this as method
func (m M20241013015640) down(db *sql.Tx) error {
	_, err := db.Exec(userDown)
	return err
}
