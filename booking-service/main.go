package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Printf("Booking service up") }).Methods("GET")

	log.Println("Listening on port 8083")
	err = http.ListenAndServe(":8083", r)
	if err != nil {
		log.Fatalf("Failed to listen on port 8083: %v", err)
	}
}
