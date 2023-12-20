package initdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB connects to the database
func initCore() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)

	var err error
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		return nil, err
	}
	// Return the db connection
	return db, nil
}

func InitDb() (*sql.DB, error) {
	var db *sql.DB
	var err error
	maxRetries := 3
	retryInterval := 4 * time.Second
	for i := 0; i < maxRetries; i++ {
		db, err = initCore()
		if err == nil {
			return db, nil
		}
		log.Printf("Failed to connect to database, attempt %d/%d: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}
	return nil, err
}
