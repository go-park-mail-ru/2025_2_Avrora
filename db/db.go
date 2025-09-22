package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	// Create tables if not exist
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Failed to create users table: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	log.Println("Database connected successfully")
}