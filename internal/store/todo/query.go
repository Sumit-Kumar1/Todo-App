package todostore

import (
	"time"

	"github.com/google/uuid"
)

const (
	insertQuery = "INSERT INTO tasks (task_id, user_id, task_title, done_status, added_at) VALUES ("
	updateQuery = "UPDATE tasks SET "
	getAll      = "SELECT task_id, user_id, task_title, done_status, added_at, modified_at from tasks WHERE user_id=?"
)

func genInsertQuery(id, title string, userID uuid.UUID, ts time.Time) (query string, values []any) {
	query = insertQuery + "?, ?, ?, ?, ?);"

	values = []any{id, userID, title, 0, ts}

	return query, values
}

func genUpdateQuery(id, title string, userID uuid.UUID, ts time.Time) (query string, vals []any) {
	query = updateQuery
	query += "task_title=?, done_status=?, modified_at=? WHERE task_id=? AND user_id=?;"

	vals = []any{title, 0, ts, id, userID}

	return
}
