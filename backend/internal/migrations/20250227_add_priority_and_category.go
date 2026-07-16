package migrations

import "database/sql"

const (
	addPriorityAndCategory = `
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS priority VARCHAR(20) DEFAULT 'MEDIUM' NOT NULL;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS category VARCHAR(100) DEFAULT '' NOT NULL;`

	dropPriorityAndCategory = `
ALTER TABLE tasks DROP COLUMN IF EXISTS priority;
ALTER TABLE tasks DROP COLUMN IF EXISTS category;`
)

type M20250227 string

// nolint:revive // unused but need this as method
func (m M20250227) up(db *sql.Tx) error {
	_, err := db.Exec(addPriorityAndCategory)
	return err
}

// nolint:revive // unused but need this as method
func (m M20250227) down(db *sql.Tx) error {
	_, err := db.Exec(dropPriorityAndCategory)
	return err
}
