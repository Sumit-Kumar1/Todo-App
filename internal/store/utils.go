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
