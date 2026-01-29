package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open database
	db, err := sql.Open("postgres", connectionString)
	// check apakah ada error saat membuka koneksi database
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	// check apakah ada error saat ping database
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (optional tapi recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Print("Database connected successfully")
	return db, nil
}
