package store

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_genInsertQuery(t *testing.T) {
	ts := time.Now()
	uid := uuid.New()
	query := "INSERT INTO tasks (task_id, task_title, done_status, added_at) VALUES (?, ?, ?, ?);"
	tests := []struct {
		name       string
		id         string
		title      string
		ts         time.Time
		wantQuery  string
		wantValues []any
	}{
		{name: "valid case", id: "abcde", title: "Dog Walk", ts: ts, wantQuery: query, wantValues: []any{"abcde", "Dog Walk", 0, ts}},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotValues := genInsertQuery(tt.id, tt.title, uid, tt.ts)

			assert.Equalf(t, tt.wantQuery, gotQuery, "TEST[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.wantValues, gotValues, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func Test_genUpdateQuery(t *testing.T) {
	ts := time.Now()
	query := "UPDATE tasks SET task_title=?, done_status=?, modified_at=? WHERE task_id=?;"
	uid := uuid.New()
	tests := []struct {
		name      string
		id        string
		title     string
		ts        time.Time
		wantQuery string
		wantVals  []any
	}{
		{name: "valid case", id: "abcde", title: "Dog Walk", ts: ts, wantQuery: query, wantVals: []any{"Dog Walk", 0, ts, "abcde"}},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotVals := genUpdateQuery(tt.id, tt.title, uid, tt.ts)

			assert.Equalf(t, tt.wantQuery, gotQuery, "TEST[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.wantVals, gotVals, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}
