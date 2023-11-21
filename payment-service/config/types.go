package config

import "fmt"

type PaymReq struct {
	BookingID int     `json:"bookingId"`
	Amount    float64 `json:"amount"`
	Curr      string  `json:"currency"`
	CardToken string  `json:"cardToken"`
}

func (p *PaymReq) Validate() error {
	if p.BookingID <= 0 {
		return fmt.Errorf("invalid booking ID")
	}
	if p.Amount <= 0.0 {
		return fmt.Errorf("invalid amount given")
	}
	if p.Curr == "" {
		return fmt.Errorf("invalid currency")
	}
	if p.CardToken == "" {
		return fmt.Errorf("invalid card token")
	}
	return nil
}
