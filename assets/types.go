package assets

import (
	"github.com/jackc/pgtype"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// LoginRequest represents the JSON structure for a login request
type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Property struct {
	ID           int     `json:"id"`
	Address      string  `json:"address"`
	ContactName  string  `json:"contact_name"`
	ContactPhone string  `json:"contact_phone"`
	ContactEmail string  `json:"contact_email"`
	Tags         string  `json:"tags"`
	Rating       float64 `json:"rating"`
	Type         string  `json:"type"`
}

type Room struct {
	ID                   int            `json:"id"`
	RoomNumber           int            `json:"room_nr"`
	HotelID              int            `json:"hotel"`
	Description          string         `json:"description"`
	Count                int            `json:"count"`
	Price                float64        `json:"price"`
	AvailabilityInterval pgtype.Tsrange `json:"availability_interval"`
}
