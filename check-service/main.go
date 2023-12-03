package main

import (
	"context"
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

func main() {
	// Initialize database with retry
	db, err := initdb.InitDb()
	if err != nil {
		log.Fatalf("could not connect to database after 3 attempts %v", err)
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
