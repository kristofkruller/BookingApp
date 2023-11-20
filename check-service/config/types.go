package config

import (
	"github.com/jackc/pgtype"
)

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
