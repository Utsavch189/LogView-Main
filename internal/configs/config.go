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

func Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		GetEnv("DB_USER"),
		GetEnv("DB_PASSWORD"),
		GetEnv("DB_HOST"),
		GetEnv("DB_PORT"),
		GetEnv("DB_SCHEMA"),
	)

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("DB open error (try %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Connected to DB successfully!")
			return db, nil
		}

		log.Printf("DB ping error (try %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("DB connect issue after retries: %v", err)
}
