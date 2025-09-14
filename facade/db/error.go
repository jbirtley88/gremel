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

func (db *ErrorGremelDB) GetSchema(tableName string) (data.Row, error) {
	return data.Row{}, db.underlyingError
}

func (db *ErrorGremelDB) DropSchema(tableName string) error {
	return db.underlyingError
}

func (db *ErrorGremelDB) GetTables() ([]string, error) {
	return nil, db.underlyingError
}

func (db *ErrorGremelDB) Mount(tableName string, source string) error {
	return db.underlyingError
}

func (db *ErrorGremelDB) GetMount(tableName string) (data.Row, error) {
	return data.Row{}, db.underlyingError
}

func (db *ErrorGremelDB) InsertRows(tableName string, rows []data.Row) error {
	return db.underlyingError
}

func (db *ErrorGremelDB) Close() error {
	return db.underlyingError
}

func (db *ErrorGremelDB) Query(sqlQuery string) ([]data.Row, []string, error) {
	return nil, nil, db.underlyingError
}
