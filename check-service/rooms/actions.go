package rooms

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
	"github.com/kristofkruller/BookingApp/check-service/config"
	"github.com/kristofkruller/BookingApp/libs/helpers"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
		return
	}

	// RoomId
	params := mux.Vars(r)
	roomIDStr := params["id"]
	// Conv
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		log.Printf("Invalid room ID format: %s", roomIDStr)
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	// QRY
	log.Printf("Fetching details for room ID: %d", roomID)
	rm := &config.Room{}
	err = db.QueryRow(`
			SELECT id, room_nr, hotel, description, count, price, availability_interval 
			FROM rooms 
			WHERE id = $1
		`, roomID).Scan(
		&rm.ID, &rm.RoomNumber, &rm.HotelID, &rm.Description, &rm.Count, &rm.Price, &rm.AvailabilityInterval,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Room ID not found: %d", roomID)
			http.Error(w, "Room ID not found", http.StatusNotFound)
		} else {
			log.Printf("Error querying database: %v", err)
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully retrieved details for room ID: %d", roomID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rm)
}

func ListRooms(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
		return
	}

	// FILTER
	var filter config.RoomFilter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Incoming filter data: %+v", filter)

	// TIME HANDL Reqperiod
	var rp *pgtype.Tsrange
	if filter.AvaStart != nil && filter.AvaEnd != nil {
		log.Printf("Ava start: %s", *filter.AvaStart)
		log.Printf("Ava end: %s", *filter.AvaEnd)

		start, err := helpers.ParseTime(*filter.AvaStart)
		if err != nil {
			log.Printf("Error parsing start date: %v", err)
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
		end, err := helpers.ParseTime(*filter.AvaEnd)
		if err != nil {
			log.Printf("Error parsing end date: %v", err)
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
		if !helpers.IsValidDateRange(start, end) {
			http.Error(w, "End date must be after start date", http.StatusBadRequest)
			return
		}

		rp = &pgtype.Tsrange{
			Lower:     pgtype.Timestamp{Time: start, Status: pgtype.Present},
			Upper:     pgtype.Timestamp{Time: end, Status: pgtype.Present},
			LowerType: pgtype.Inclusive,
			UpperType: pgtype.Exclusive,
		}

		log.Printf("Parsed availability period: [%s, %s]", start.Format(time.RFC3339), end.Format(time.RFC3339))
	}

	// QRY base
	q := []string{`
		SELECT id, room_nr, hotel, description, count, price, availability_interval 
		FROM rooms
	`}
	// QRY params
	var args []interface{}
	var condi []string

	if filter.PriceMin != nil {
		condi = append(condi, "price >= $1")
		args = append(args, filter.PriceMin)
	}
	if filter.PriceMax != nil {
		condi = append(condi, "price <= $2")
		args = append(args, filter.PriceMax)
	}
	if filter.PriceMin != nil && filter.PriceMax != nil {
		if !helpers.IsValidPriceLogic(*filter.PriceMin, *filter.PriceMax) {
			http.Error(w, "Min price must be lower or equal to max price", http.StatusBadRequest)
			return
		}
	}
	if len(condi) > 0 {
		q = append(q, "WHERE "+strings.Join(condi, " AND "))
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
	var rooms []*config.Room
	for rws.Next() {
		var room config.Room
		err := rws.Scan(&room.ID, &room.RoomNumber, &room.HotelID, &room.Description, &room.Count, &room.Price, &room.AvailabilityInterval)
		if err != nil {
			http.Error(w, "Error reading room data", http.StatusInternalServerError)
			return
		}

		log.Printf("Fetched room: %v", room)
		if rp != nil && !helpers.IsRoomAvailable(room.ID, *rp, db) {
			log.Printf("Room ID %d is not available for the requested period", room.ID)
			continue // skip
		}

		log.Printf("Room ID %d is available and added to the response", room.ID)
		rooms = append(rooms, &room)
	}
	// Handle any error encountered during iteration
	if err = rws.Err(); err != nil {
		log.Printf("Error iterating room data: %v", err)
		http.Error(w, "Error iterating rooms data", http.StatusInternalServerError)
		return
	}

	log.Println("Successfully retrieved rooms")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}
