package migrations

import (
	"github.com/sqlitecloud/sqlitecloud-go"
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
func (m M20241013015640) up(db *sqlitecloud.SQCloud) error {
	return db.Execute(userUp)
}

// nolint:revive // unused but need this as method
func (m M20241013015640) down(db *sqlitecloud.SQCloud) error {
	return db.Execute(userDown)
}
