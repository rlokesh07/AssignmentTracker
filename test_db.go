package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgresql://postgres.mztcyklvjxbyokhdzehl:R1shi5538!!@aws-0-us-east-1.pooler.supabase.com:6543/postgres"

	fmt.Println("Attempting to connect...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Open error:", err)
	}
	defer db.Close()

	fmt.Println("Pinging database...")
	err = db.Ping()
	if err != nil {
		log.Fatal("Ping error:", err)
	}

	fmt.Println("Connected successfully!")

	fmt.Println("Testing query...")
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM assignments").Scan(&count)
	if err != nil {
		log.Fatal("Query error:", err)
	}

	fmt.Printf("Found %d assignments\n", count)
}
