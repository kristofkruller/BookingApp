package reserv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
	"github.com/kristofkruller/BookingApp/booking-service/config"
	"github.com/kristofkruller/BookingApp/booking-service/helpers"
)

func LetsBook(w http.ResponseWriter, r *http.Request) {
	st := time.Now()

	var br config.BookingReq
	err := json.NewDecoder(r.Body).Decode(&br)
	if err != nil {
		log.Printf("CreateBooking: Error decoding request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if err := br.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TIME HANDL Bookingperiod
	var bp *pgtype.Tsrange
	log.Printf("Booking start: %s", br.StartDate)
	log.Printf("Booking end: %s", br.EndDate)

	start, err := helpers.ParseTime(br.StartDate)
	if err != nil {
		log.Printf("Error parsing start date: %v", err)
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}
	end, err := helpers.ParseTime(br.EndDate)
	if err != nil {
		log.Printf("Error parsing end date: %v", err)
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}
	if !helpers.IsValidDateRange(start, end) {
		http.Error(w, "End date must be after start date", http.StatusBadRequest)
		return
	}

	bp = &pgtype.Tsrange{
		Lower:     pgtype.Timestamp{Time: start, Status: pgtype.Present},
		Upper:     pgtype.Timestamp{Time: end, Status: pgtype.Present},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
	}

	log.Printf("Parsed availability period: [%s, %s]", start.Format(time.RFC3339), end.Format(time.RFC3339))

	// AVA
	if bp != nil && !isRoomAvailable(br.RoomID, *bp) {
		log.Printf("Room ID %d is not available for the requested period", br.RoomID)
		http.Error(w, "Room is not available for the selected dates", http.StatusConflict)
		return
	}

	// INSERT INTO DB
	_, err = db.Exec(`
		INSERT INTO reserv (
			userId, propertyId, roomId, cost, reserv_interval
		) VALUES ($1, $2, $3, $4, tsrange($5, $6, '[]'))`,
		br.UserID, br.PropertyID, br.RoomID, br.Cost, br.StartDate, br.EndDate,
	)
	if err != nil {
		log.Printf("CreateBooking: Error creating booking: %v", err)
		http.Error(w, "Error creating booking", http.StatusInternalServerError)
		return
	}

	log.Printf("CreateBooking: Booking created successfully in %v", time.Since(st))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Booking created successfully")
}

func DontBook(w http.ResponseWriter, r *http.Request) {
	st := time.Now()

	vars := mux.Vars(r)
	bId := vars["bookingId"]

	// VALIDATE
	if _, err := strconv.Atoi(bId); err != nil {
		log.Printf("CancelBooking: Invalid booking ID: %v", err)
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// DELETE FROM DB
	res, err := db.Exec(`
			DELETE FROM reserv 
			WHERE id = $1
		`, bId)
	if err != nil {
		log.Printf("CancelBooking: Error canceling booking: %v", err)
		http.Error(w, "Error canceling booking", http.StatusInternalServerError)
		return
	}

	// CHECK ROWS AFFECTED OR NOT
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("CancelBooking: Error getting rows affected: %v", err)
		http.Error(w, "Error during cancellation", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No booking found with the given ID", http.StatusNotFound)
		return
	}

	log.Printf("CancelBooking: Booking canceled successfully in %v", time.Since(st))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Booking canceled successfully")
}

// HELPERS
func isRoomAvailable(roomID int, period pgtype.Tsrange) bool {
	var count int
	var lower, upper pgtype.Timestamp
	lower.Set(period.Lower.Time)
	upper.Set(period.Upper.Time)

	tsrange := pgtype.Tsrange{
		Lower:     lower,
		Upper:     upper,
		LowerType: period.LowerType,
		UpperType: period.UpperType,
		Status:    pgtype.Present,
	}

	err := db.QueryRow(`
			SELECT COUNT(*) 
			FROM reserv 
			WHERE roomId = $1 
			AND reserv_interval && $2
	`, roomID, tsrange).Scan(&count)

	if err != nil {
		log.Printf("Error in isRoomAvailable: %v", err)
		return false
	}

	return count == 0
}
