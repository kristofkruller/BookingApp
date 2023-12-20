package reserv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
	"github.com/kristofkruller/BookingApp/booking-service/config"
	"github.com/kristofkruller/BookingApp/libs/helpers"
)

var db *sql.DB

func BookingsOfHTTP(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
		return
	}

	// UserId
	params := mux.Vars(r)
	uIdStr := params["uId"]
	// Conv
	uId, err := strconv.Atoi(uIdStr)
	if err != nil {
		log.Printf("Invalid user ID format: %s", uIdStr)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// FILTER
	var filter config.BookingFilter
	err = json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// QRY base
	log.Printf("Fetching details for user ID: %d", uId)
	q := []string{`
		SELECT id, userId, propertyId, roomId, cost, reserv_interval, creation_date 
		FROM reserv 
		WHERE userId = $1
	`}

	// QRY params
	var args []interface{}
	args = append(args, uId)
	var condi []string

	argCount := 2 // $1 for uId

	// TIME HANDL crDate
	if filter.CreationDate != nil {
		creationDate, err := helpers.ParseTime(*filter.CreationDate)
		if err != nil {
			log.Printf("Error parsing creation date: %v", err)
			http.Error(w, "Invalid creation date format", http.StatusBadRequest)
			return
		}
		condi = append(condi, "DATE(creation_date) = $"+strconv.Itoa(argCount))
		args = append(args, creationDate)
		argCount++
	}

	// Start, End dates
	if filter.StartDate != nil && filter.EndDate != nil {
		log.Printf("StartDate: %s", *filter.StartDate)
		log.Printf("EndDate: %s", *filter.EndDate)

		start, err := helpers.ParseTime(*filter.StartDate)
		if err != nil {
			log.Printf("Error parsing start date: %v", err)
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
		end, err := helpers.ParseTime(*filter.EndDate)
		if err != nil {
			log.Printf("Error parsing end date: %v", err)
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
		if !helpers.IsValidDateRange(start, end) {
			http.Error(w, "End date must be after start date", http.StatusBadRequest)
			return
		}

		condi = append(condi, "reserv_interval && tsrange($"+strconv.Itoa(argCount)+", $"+strconv.Itoa(argCount+1)+", '[]')")
		args = append(args, start, end)
		argCount += 2
	}

	if filter.PriceMin != nil && filter.PriceMax != nil {
		if !helpers.IsValidPriceLogic(*filter.PriceMin, *filter.PriceMax) {
			http.Error(w, "Invalid price range", http.StatusBadRequest)
			return
		}
		condi = append(condi, "price >= $"+strconv.Itoa(argCount))
		args = append(args, *filter.PriceMin)
		argCount++

		condi = append(condi, "price <= $"+strconv.Itoa(argCount))
		args = append(args, *filter.PriceMax)
		argCount++
	}

	if len(condi) > 0 {
		q = append(q, " AND "+strings.Join(condi, " AND "))
	}
	// QRY filtered
	qry := strings.Join(q, " ")

	// Dereference args for logging
	var logArgs []interface{}
	for _, arg := range args {
		if argPointer, ok := arg.(*float64); ok && argPointer != nil {
			logArgs = append(logArgs, *argPointer)
		} else {
			logArgs = append(logArgs, arg)
		}
	}

	log.Printf("Executing query: %s with args: %v", qry, logArgs)

	// QRY fire - rows as result
	rws, err := db.Query(qry, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Error at building up query", http.StatusInternalServerError)
		return
	}

	defer rws.Close()

	// PROCESSING AND BUILDING UP RES
	var bookings []*config.Booking

	for rws.Next() {
		var booking config.Booking
		if err := rws.Scan(&booking.ID, &booking.UserID, &booking.PropertyID, &booking.RoomID, &booking.Cost, &booking.ReservInterval, &booking.CreationDate); err != nil {
			log.Printf("Error reading booking data: %v", err)
			http.Error(w, "Error reading booking data", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, &booking)
	}

	// Check for errors after iteration
	if err = rws.Err(); err != nil {
		log.Printf("Error iterating booking data: %v", err)
		http.Error(w, "Error iterating booking data", http.StatusInternalServerError)
		return
	}

	// RETURN
	// as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Printf("Error encoding bookings to JSON: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// func printBookingsToConsole(bookings []*config.Booking) {
// 	log.Println("Booking History:")
// 	for _, booking := range bookings {
// 		log.Printf("ID: %d, UserID: %d, PropertyID: %d, RoomID: %d, Cost: %.2f, Interval: %v, Creation Date: %v\n",
// 			booking.ID, booking.UserID, booking.PropertyID, booking.RoomID, booking.Cost, booking.ReservInterval, booking.CreationDate)
// 	}
// }

func LetsBookHTTP(w http.ResponseWriter, r *http.Request) {
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
	if bp != nil && !helpers.IsRoomAvailable(br.RoomID, *bp, db) {
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
		log.Printf("Error creating booking: %v", err)
		http.Error(w, "Error creating booking", http.StatusInternalServerError)
		return
	}

	// all booking of user refreshed directly from db for validation
	var cb []config.Booking
	cbrows, err := db.Query(`
		SELECT *
		FROM reserv
		WHERE userId = $1`,
		&br.UserID,
	)
	if err != nil {
		log.Printf("Error validating bookings of user: %v", err)
		http.Error(w, "Error validating booking", http.StatusInternalServerError)
		return
	}
	defer cbrows.Close()
	// Process rows and append to bookings slice
	for cbrows.Next() {
		var b config.Booking
		if err := cbrows.Scan(&b.ID, &b.UserID, &b.PropertyID, &b.RoomID, &b.Cost, &b.ReservInterval, &b.CreationDate); err != nil {
			http.Error(w, "Error scanning booking row: %w", http.StatusInternalServerError)
			return
		}
		cb = append(cb, b)
	}
	// Check for errors after iterating through rows
	if err := cbrows.Err(); err != nil {
		http.Error(w, "Error iterating booking rows: %w", http.StatusInternalServerError)
	}

	log.Printf("Booking created successfully in %v", time.Since(st))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(cb); err != nil {
		log.Printf("Error encoding bookings to JSON: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func DontBookHTTP(w http.ResponseWriter, r *http.Request) {
	st := time.Now()

	vars := mux.Vars(r)
	bId := vars["bookingId"]

	// VALIDATE
	if _, err := strconv.Atoi(bId); err != nil {
		log.Printf("Invalid booking ID az cancellation of booking: %v", err)
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// DELETE FROM DB
	res, err := db.Exec(`
			DELETE FROM reserv 
			WHERE id = $1
		`, bId)
	if err != nil {
		log.Printf("Error canceling booking: %v", err)
		http.Error(w, "Error canceling booking", http.StatusInternalServerError)
		return
	}

	// CHECK ROWS AFFECTED OR NOT
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Error during cancellation", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No booking found with the given ID", http.StatusNotFound)
		return
	}

	log.Printf("Booking canceled successfully in %v", time.Since(st))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Booking canceled successfully")
}
