package config

import (
	"time"

	"github.com/jackc/pgtype"
)

type Booking struct {
	ID             int            `json:"id"`
	UserID         int            `json:"userId"`
	PropertyID     int            `json:"propertyId"`
	RoomID         int            `json:"roomId"`
	Cost           float64        `json:"cost"`
	ReservInterval pgtype.Tsrange `json:"reserv_interval"`
	CreationDate   time.Time      `json:"creation_date"`
}

type BookingFilter struct {
	PriceMin     *float64 `json:"min_price,omitempty"`
	PriceMax     *float64 `json:"max_price,omitempty"`
	CreationDate *string  `json:"creation_date,omitempty"`
	StartDate    *string  `json:"start_date,omitempty"`
	EndDate      *string  `json:"end_date,omitempty"`
}
