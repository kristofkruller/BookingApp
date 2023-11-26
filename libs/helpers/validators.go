package helpers

import (
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgtype"
)

func IsValidDateRange(start, end time.Time) bool {
	return start.Before(end) || start.Equal(end)
}

func IsValidPriceLogic(min, max float64) bool {
	return min <= max
}

func IsRoomAvailable(roomID int, period pgtype.Tsrange, db *sql.DB) bool {
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
