package config

import (
	"fmt"
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

type BookingReq struct {
	UserID     int     `json:"userId"`
	PropertyID int     `json:"propertyId"`
	RoomID     int     `json:"roomId"`
	Cost       float64 `json:"cost"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
}

func (b *BookingReq) Validate() error {
	if b.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if b.PropertyID <= 0 {
		return fmt.Errorf("invalid property ID")
	}
	if b.RoomID <= 0 {
		return fmt.Errorf("invalid room ID")
	}
	if b.Cost <= 0.0 {
		return fmt.Errorf("invalid cost")
	}
	if b.StartDate == "" {
		return fmt.Errorf("start date is required")
	}
	if b.EndDate == "" {
		return fmt.Errorf("end date is required")
	}
	return nil
}
