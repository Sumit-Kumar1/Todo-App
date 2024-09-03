package store

import "database/sql"

func rollback(tx *sql.Tx, err error) error {
	if err == nil {
		return nil
	}

	if rlErr := tx.Rollback(); rlErr != nil {
		return rlErr
	}

	return err
}

func runMigration(db *sql.DB) error {
	const (
		upTaskTable = `DROP TABLE IF EXISTS tasks; CREATE TABLE IF NOT EXISTS tasks(task_id TEXT PRIMARY KEY, user_id TEXT NOT NULL,
task_title TEXT NOT NULL, done_status BOOLEAN NOT NULL CHECK (done_status IN (0, 1)),
added_at DATETIME NOT NULL, modified_at DATETIME);`
		upUserTable = `DROP TABLE IF EXISTS users; CREATE TABLE IF NOT EXISTS users(user_id TEXT NOT NULL PRIMARY KEY, name TEXT NOT NULL,
email TEXT NOT NULL UNIQUE CHECK (email LIKE '%'), password TEXT NOT NULL);`
		upSessionTable = `DROP TABLE IF EXISTS sessions; CREATE TABLE IF NOT EXISTS sessions(id TEXT PRIMARY KEY, user_id TEXT NOT NULL UNIQUE,
token TEXT NOT NULL, expiry DATETIME NOT NULL);`
	)

	if _, err := db.Exec(upTaskTable); err != nil {
		return err
	}

	if _, err := db.Exec(upUserTable); err != nil {
		return err
	}

	if _, err := db.Exec(upSessionTable); err != nil {
		return err
	}

	return nil
}
