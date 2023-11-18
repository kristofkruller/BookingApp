package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kristofkruller/BookingApp/auth-service/config"
	"github.com/kristofkruller/BookingApp/auth-service/users"
)

func main() {
	//DB
	err := config.InitDB("postgres://user:password@localhost/BookingAppDb?sslmode=disable")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	//ROUTER
	r := mux.NewRouter()

	r.HandleFunc("/login", users.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", users.LogoutHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", r))
}
