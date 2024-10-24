package migrations

import (
	"database/sql"
)

const (
	userDown = "DROP TABLE IF EXISTS users;"
	userUp   = `CREATE TABLE IF NOT EXISTS users(
    user_id TEXT NOT NULL PRIMARY KEY, 
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE CHECK (email LIKE '%'), 
    password TEXT NOT NULL);`
)

type M20241013015640 string

func (m M20241013015640) Up(db *sql.Tx) error {
	return run(db, userUp, m)
}

func (m M20241013015640) Down(db *sql.Tx) error {
	return run(db, userDown, m)
}
