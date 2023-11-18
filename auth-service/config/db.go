package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB connects to the database
func InitDB(dbConnString string) error {
	var err error
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		return err
	}
	return db.Ping()
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}
