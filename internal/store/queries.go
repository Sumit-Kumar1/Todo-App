package store

import (
	"time"
)

const (
	insertQuery   = "INSERT INTO tasks (task_id, task_title, done_status, added_at) VALUES ("
	updateQuery   = "UPDATE tasks SET "
	getAll        = "SELECT task_id, task_title, done_status, added_at, modified_at from tasks"
	registerQuery = "INSERT INTO users(user_id, name, email, password) VALUES (?,?,?,?);"
	getUser       = "SELECT * FROM users WHERE email = ?;"
)

func genInsertQuery(id, title string, ts time.Time) (query string, values []any) {
	query = insertQuery + "?, ?, ?, ?);"

	values = []any{id, title, 0, ts}

	return query, values
}

func genUpdateQuery(id, title string, ts time.Time) (query string, vals []any) {
	query = updateQuery
	query += "task_title=?, done_status=?, modified_at=? WHERE task_id=?;"

	vals = []any{title, 0, ts, id}

	return
}
