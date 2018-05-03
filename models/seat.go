package models

// Seat Model
type Seat struct {
	SeatID     int    `json:"seat_id"`
	Number     string `json:"number"`
	QRCodeUrl  string `json:"qr_code_url"`
	MerchantID int    `json:"merchant_id"`
}
