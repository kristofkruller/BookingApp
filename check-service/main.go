package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kristofkruller/BookingApp/check-service/rooms"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Check-service up")
	}).Methods("GET")
	r.HandleFunc("/room/{id}", rooms.GetRoom).Methods("GET")
	// r.HandleFunc("/rooms", rooms.ListRooms).Methods("GET")

	log.Println("Listening on port 8082")
	err := http.ListenAndServe(":8082", r)
	if err != nil {
		log.Fatalf("Failed to listen on port 8082: %v", err)
	}
}
