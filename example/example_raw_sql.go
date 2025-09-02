package example

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func RawSQL() {
	// Create two completely separate named in-memory SQLite databases
	// Using file: syntax with cache=shared ensures each named database is separate
	category_db, err := sql.Open("sqlite3", "file:category_db?mode=memory&cache=shared")
	if err != nil {
		log.Fatalf("Failed to open category_db: %v", err)
	}
	defer category_db.Close()

	description_db, err := sql.Open("sqlite3", "file:description_db?mode=memory&cache=shared")
	if err != nil {
		log.Fatalf("Failed to open description_db: %v", err)
	}
	defer description_db.Close()

	// Test the connections
	if err := category_db.Ping(); err != nil {
		log.Fatalf("Failed to ping db1: %v", err)
	}

	if err := description_db.Ping(); err != nil {
		log.Fatalf("Failed to ping db2: %v", err)
	}

	// Create different tables in each database to prove they're separate
	_, err = category_db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)`)
	if err != nil {
		log.Fatalf("Failed to create table in db1: %v", err)
	}

	_, err = description_db.Exec(`CREATE TABLE products (id INTEGER PRIMARY KEY, description TEXT, price REAL)`)
	if err != nil {
		log.Fatalf("Failed to create table in db2: %v", err)
	}

	// Insert data into db1
	_, err = category_db.Exec(`INSERT INTO users (name) VALUES ('Alice'), ('Bob')`)
	if err != nil {
		log.Fatalf("Failed to insert into db1: %v", err)
	}

	// Insert data into db2
	_, err = description_db.Exec(`INSERT INTO products (description, price) VALUES ('Widget A', 19.99), ('Widget B', 29.99)`)
	if err != nil {
		log.Fatalf("Failed to insert into db2: %v", err)
	}

	// Query data from db1
	rows1, err := category_db.Query(`SELECT id, name FROM users`)
	if err != nil {
		log.Fatalf("Failed to query db1: %v", err)
	}
	defer rows1.Close()

	fmt.Println("Data from db1 (users):")
	for rows1.Next() {
		var id int
		var name string
		if err := rows1.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  ID: %d, Name: %s\n", id, name)
	}

	// Query data from db2
	rows2, err := description_db.Query(`SELECT id, description, price FROM products`)
	if err != nil {
		log.Fatalf("Failed to query db2: %v", err)
	}
	defer rows2.Close()

	fmt.Println("\nData from db2 (products):")
	for rows2.Next() {
		var id int
		var description string
		var price float64
		if err := rows2.Scan(&id, &description, &price); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  ID: %d, Description: %s, Price: $%.2f\n", id, description, price)
	}

	fmt.Println("\n✅ Successfully demonstrated two completely separate named in-memory SQLite databases!")
}
