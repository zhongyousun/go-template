package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := os.Getenv("POSTGRES_CONN")
	var err error

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	if err := DB.Ping(); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
}
