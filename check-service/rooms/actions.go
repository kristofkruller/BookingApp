package rooms

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kristofkruller/BookingApp/assets"
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

	// roomId
	params := mux.Vars(r)
	roomIDStr := params["id"]
	// conv
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	// qry
	rm := &assets.Room{}
	err = db.QueryRow("SELECT id, room_nr, hotel, description, count, price, availability_interval FROM rooms WHERE id = $1", roomID).Scan(
		&rm.ID, &rm.RoomNumber, &rm.HotelID, &rm.Description, &rm.Count, &rm.Price, &rm.AvailabilityInterval,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(rm)
}
