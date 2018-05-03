package models

import (
	"time"
)

// Order Model
type Order struct {
	OrderID      int       `json:"order_id"`
	Status       string    `json:"status"`
	SeatID       int       `json:"seat_id"`
	CustomerID   int       `json:"customer_id"`
	MerchantID   int       `json:"merchant_id"`
	OrderTime    time.Time `json:"order_time"`
	CompleteTime time.Time `json:"complete_time"`
}
