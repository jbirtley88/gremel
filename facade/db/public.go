package db

import "github.com/jbirtley88/gremel/data"

type GremelDB interface {
	CreateSchema(tableName string, row data.Row) error
	DropSchema(tableName string) error
	InsertRows(tableName string, rows []data.Row) error
	Close() error
}

func NewGremelDB(name string) GremelDB {
	// Only sqlite for now
	return NewNamedSQLiteGremelDB(name)
}

/*
-- Create the accounts table
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    mac_address VARCHAR(17) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    percent INTEGER NOT NULL CHECK (percent >= 0 AND percent <= 100)
);

-- Create indexes for better query performance
CREATE INDEX idx_accounts_username ON accounts(username);
CREATE INDEX idx_accounts_email ON accounts(email);
CREATE INDEX idx_accounts_mac_address ON accounts(mac_address);
CREATE INDEX idx_accounts_percent ON accounts(percent);
*/
