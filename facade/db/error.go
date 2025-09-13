package db

import "github.com/jbirtley88/gremel/data"

type ErrorGremelDB struct {
	underlyingError error
}

func NewErrorGremelDB(err error) GremelDB {
	return &ErrorGremelDB{
		underlyingError: err,
	}
}

func (db *ErrorGremelDB) CreateSchema(tableName string, row data.Row) error {
	return db.underlyingError
}

func (db *ErrorGremelDB) DropSchema(tableName string) error {
	return db.underlyingError
}

func (db *ErrorGremelDB) InsertRows(tableName string, rows []data.Row) error {
	return db.underlyingError
}

func (db *ErrorGremelDB) Close() error {
	return db.underlyingError
}
