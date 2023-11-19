package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/kristofkruller/BookingApp/auth-service/users"
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB connects to the database
func InitDB() (*sql.DB, error) {
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

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// SetDB allows other packages to set the database connection
func SetDB(database *sql.DB) {
	users.SetDB(database)
}
