package userstore

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func initTest() (*Store, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Store{DB: db, Log: logger}, mock, nil
}

func TestStore_Logout(t *testing.T) {
	store, mock, err := initTest()
	if err != nil {
		store.Log.Error(err.Error())
		return
	}

	tests := []struct {
		name    string
		token   *uuid.UUID
		wantErr error
	}{
		{"valid case"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantErr, store.Logout(context.Background(), tt.token))
		})
	}
}
