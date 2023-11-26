package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kristofkruller/BookingApp/check-service/rooms"
	"github.com/kristofkruller/BookingApp/libs/initdb"
)

const (
	maxRetries    = 5               // Maximum number of retries
	retryInterval = 5 * time.Second // Time to wait between retries
)

func main() {
	var err error

	// Initialize database with retry
	var db *sql.DB
	for i := 0; i < maxRetries; i++ {
		db, err = initdb.InitDB()
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database, attempt %d/%d: %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
	}

	//DB conn for rooms pkg
	rooms.SetDB(db)

	//ROUTER
	r := mux.NewRouter()

	r.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Check-service up")
	}).Methods("GET")
	r.HandleFunc("/room/{id}", rooms.GetRoom).Methods("GET")
	r.HandleFunc("/rooms", rooms.ListRooms).Methods("POST")

	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	// Start HTTP Server in a Goroutine
	go func() {
		log.Println("Listening on port 8082")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
