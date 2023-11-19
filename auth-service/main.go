package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kristofkruller/BookingApp/auth-service/config"
	"github.com/kristofkruller/BookingApp/auth-service/users"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//DB
	err = config.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	//ROUTER
	r := mux.NewRouter()

	r.HandleFunc("/login", users.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", users.LogoutHandler).Methods("POST")

	log.Println("Listening on port 8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatalf("Failed to listen on port 8081: %v", err)
	}
}
