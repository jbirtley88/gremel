package db

import (
	"github.com/jbirtley88/gremel/data"
)

type GremelDB interface {
	CreateSchema(tableName string, row data.Row) error
	DropSchema(tableName string) error
	// Support for the '.schema' command
	GetSchema(tableName string) (data.Row, error)
	// Support for the '.tables' command
	GetTables() ([]string, error)

	// Register the mount point for this table, to support the '.mount' command
	Mount(tableName string, source string) error
	// Get the mount point for this table, to support the '.mount' command
	GetMount(tableName string) (data.Row, error)
	InsertRows(tableName string, rows []data.Row) error
	Query(sqlQuery string) ([]data.Row, []string, error)
	Close() error
}
