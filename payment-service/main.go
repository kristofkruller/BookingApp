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
	"github.com/kristofkruller/BookingApp/payment-service/processor"
)

func main() {
	//ROUTER
	r := mux.NewRouter()

	r.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Payment-service up")
	}).Methods("GET")
	r.HandleFunc("/pay/{bookingId}", processor.Pay).Methods("POST")

	srv := &http.Server{
		Addr:    ":8084",
		Handler: r,
	}

	// Start HTTP Server in a Goroutine
	go func() {
		log.Println("Listening on port 8084")
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
