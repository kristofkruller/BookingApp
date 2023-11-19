package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbConStr := os.Getenv("DB_CONNECTION_SEED")
	aPass := os.Getenv("DB_PASSWORD")

	db, err := sql.Open("postgres", dbConStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Ensure the database is ready
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Hash the admin password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(aPass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	// Insert or update the admin user
	_, err = db.Exec(`INSERT INTO users (name, email, password) VALUES ('admin', 'admin@example.com', $1)
                      ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password`, string(hashedPassword))
	if err != nil {
		log.Fatalf("Error seeding admin user: %v", err)
	}

	fmt.Println("Database seeding completed successfully")
}
