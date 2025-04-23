package configs

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
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

func Connect() (*sql.DB, error) {
	dbPath := "logview_data/logview.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening SQLite connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging SQLite: %v", err)
	}

	return db, nil
}
