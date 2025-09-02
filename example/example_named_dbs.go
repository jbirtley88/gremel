package example

import (
	"fmt"
	"log"

	"github.com/jbirtley88/gremel/data"
	"github.com/jbirtley88/gremel/facade/db"
)

func NamedDBs() {
	// Create two completely separate named in-memory SQLite databases
	category_db := db.NewNamedSQLiteGremelDB("category_db")
	description_db := db.NewNamedSQLiteGremelDB("description_db")

	// Make sure to close both databases when done
	defer func() {
		if err := category_db.Close(); err != nil {
			log.Printf("Error closing category_db: %v", err)
		}
	}()
	defer func() {
		if err := description_db.Close(); err != nil {
			log.Printf("Error closing description_db: %v", err)
		}
	}()

	// Create sample data for both databases
	sampleRow1 := data.Row{
		"id":       1,
		"name":     "Category DB Record",
		"category": "test",
	}

	sampleRow2 := data.Row{
		"id":          1,
		"description": "Description DB Record",
		"value":       42.5,
	}

	// Create schemas for both databases with different structures
	fmt.Println("Creating schema for category_db...")
	err := category_db.CreateSchema("table1", sampleRow1)
	if err != nil {
		log.Fatalf("Failed to create schema for category_db: %v", err)
	}

	fmt.Println("Creating schema for description_db...")
	err = description_db.CreateSchema("table2", sampleRow2)
	if err != nil {
		log.Fatalf("Failed to create schema for description_db: %v", err)
	}

	// Insert data into both databases
	fmt.Println("Inserting data into category_db...")
	rows1 := []data.Row{
		{
			"id":       1,
			"name":     "First Record",
			"category": "A",
		},
		{
			"id":       2,
			"name":     "Second Record",
			"category": "B",
		},
	}
	err = category_db.InsertRows("table1", rows1)
	if err != nil {
		log.Fatalf("Failed to insert into category_db: %v", err)
	}

	fmt.Println("Inserting data into description_db...")
	rows2 := []data.Row{
		{
			"id":          1,
			"description": "First Description",
			"value":       10.5,
		},
		{
			"id":          2,
			"description": "Second Description",
			"value":       20.7,
		},
	}
	err = description_db.InsertRows("table2", rows2)
	if err != nil {
		log.Fatalf("Failed to insert into description_db: %v", err)
	}

	fmt.Println("✅ Successfully created and populated two separate named in-memory SQLite databases!")
	fmt.Println("- category_db contains table1 with id, name, category columns")
	fmt.Println("- description_db contains table2 with id, description, value columns")
	fmt.Println("- Both databases are completely independent")
}
