package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

func main() {
	// Try different connection string formats
	connStrings := []string{
		"host=db.mztcyklvjxbyokhdzehl.supabase.co port=5432 user=postgres password=R1shi17905538~~ dbname=postgres sslmode=require",
		"postgres://postgres:R1shi17905538~~@db.mztcyklvjxbyokhdzehl.supabase.co:5432/postgres?sslmode=require",
		"host=db.mztcyklvjxbyokhdzehl.supabase.co port=5432 user=postgres password=R1shi17905538~~ dbname=postgres sslmode=disable",
	}

	for i, connStr := range connStrings {
		fmt.Printf("\n=== Trying connection string %d ===\n", i+1)
		fmt.Printf("Format: %s\n", connStr)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Failed to open connection: %v", err)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping: %v", err)
			db.Close()
			continue
		}

		fmt.Println("✓ Connection successful!")

		// Try a simple query
		var result int
		err = db.QueryRow("SELECT 1").Scan(&result)
		if err != nil {
			log.Printf("Query failed: %v", err)
		} else {
			fmt.Printf("✓ Query successful! Result: %d\n", result)
		}

		db.Close()
		return
	}

	fmt.Println("\n❌ All connection attempts failed")
}
