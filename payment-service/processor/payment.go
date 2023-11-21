package processor

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kristofkruller/BookingApp/payment-service/config"
)

func Pay(w http.ResponseWriter, r *http.Request) {
	st := time.Now()

	vars := mux.Vars(r)
	bId := vars["bookingId"]

	// VALIDATE
	if _, err := strconv.Atoi(bId); err != nil {
		log.Printf("CancelBooking: Invalid booking ID: %v", err)
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// DECODE
	var pr config.PaymReq
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		log.Printf("Error decoding payment request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := pr.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Processing payment for Booking ID: %d, Amount: %.2f %s", pr.BookingID, pr.Amount, pr.Curr)

	// Simulate payment processing
	success := mockPaymentProcessor(&pr)
	if success {
		log.Printf("Payment successful for Booking ID: %d", pr.BookingID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	} else {
		log.Printf("Payment failed for Booking ID: %d", pr.BookingID)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "failed"})
	}

	elapsedTime := time.Since(st)
	log.Printf("Processed payment for Booking ID: %d in %s", pr.BookingID, elapsedTime)
}

// HELPERS
func mockPaymentProcessor(req *config.PaymReq) bool {
	// Simulate some processing time
	time.Sleep(2 * time.Second)

	// Randomly decide
	return rand.Intn(2) == 0
}
