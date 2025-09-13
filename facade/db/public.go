package db

import (
	"github.com/jbirtley88/gremel/data"
)

type GremelDB interface {
	CreateSchema(tableName string, row data.Row) error
	DropSchema(tableName string) error
	InsertRows(tableName string, rows []data.Row) error
	Query(sqlQuery string) ([]data.Row, error)
	Close() error
}
