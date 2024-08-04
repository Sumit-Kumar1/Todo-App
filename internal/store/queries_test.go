package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_genInsertQuery(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		title      string
		ts         time.Time
		wantQuery  string
		wantValues []any
	}{
		// TODO: Add test cases.
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotValues := genInsertQuery(tt.id, tt.title, tt.ts)

			assert.Equalf(t, tt.wantQuery, gotQuery, "TEST[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.wantValues, gotValues, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}
