package configs

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func GetEnv(key string, fallback ...string) string {
	_ = godotenv.Load()

	defaultValue := ""
	if len(fallback) > 0 {
		defaultValue = fallback[0]
	}
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SQLLite connection

// func Connect() (*sql.DB, error) {
// 	dbPath := "logview_data/logview.db"

// 	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
// 	if err != nil {
// 		return nil, fmt.Errorf("error opening SQLite connection: %v", err)
// 	}

// 	if err := db.Ping(); err != nil {
// 		return nil, fmt.Errorf("error pinging SQLite: %v", err)
// 	}

// 	return db, nil
// }

// MySQL connection

var db *sql.DB

func Connect() (*sql.DB, error) {
	// If the db connection is already established, reuse it
	if db != nil {
		return db, nil
	}

	// Build the Data Source Name (DSN) from environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		GetEnv("DB_USER"),
		GetEnv("DB_PASSWORD"),
		GetEnv("DB_HOST"),
		GetEnv("DB_PORT"),
		GetEnv("DB_SCHEMA"),
	)

	// Open a new database connection
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("DB open error: %v", err)
		return nil, err
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Check if the DB is reachable by pinging it
	err = db.Ping()
	if err != nil {
		log.Printf("DB ping error: %v", err)
		return nil, err
	}

	log.Println("Connected to DB successfully!")
	return db, nil
}
