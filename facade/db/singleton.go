package db

import "sync"

var dbSingleton sync.Once
var dbInstance GremelDB

// GetGremelDB returns a singleton instance of the GremelDB
// This is a simple way to ensure that we only have one database instance in the application
// and that it is shared across all components that need it.
func GetGremelDB() GremelDB {
	dbSingleton.Do(func() {
		dbInstance = newNamedSQLiteGremelDB("gremel")
	})
	return dbInstance
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
