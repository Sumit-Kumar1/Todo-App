package migrations

import (
	"github.com/sqlitecloud/sqlitecloud-go"
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

func (m M20241013015640) Up(db *sqlitecloud.SQCloud) error {
	return run(db, userUp, m)
}

func (m M20241013015640) Down(db *sqlitecloud.SQCloud) error {
	return run(db, userDown, m)
}
