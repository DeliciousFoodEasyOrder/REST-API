package models

// Seat Model
type Seat struct {
	SeatID     int    `json:"seat_id"`
	Number     string `json:"number"`
	QRCodeURL  string `json:"qr_code_url"`
	MerchantID int    `json:"merchant_id"`
}
