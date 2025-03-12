package store

import "github.com/sqlitecloud/sqlitecloud-go"

//go:generate mockgen --source=interface.go --destination=mock_interface.go --package=store
type SqliteClouder interface {
	Select(SQL string) (*sqlitecloud.Result, error)
	Execute(SQL string) error
}
