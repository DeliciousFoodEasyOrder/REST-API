package models

import (
	"time"
)

// Merchant Model
type Merchant struct {
	MerchantID int    `json:"merchant_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
}

// Food Model
type Food struct {
	FoodID      int     `json:"food_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	MerchantID  int     `json:"merchant_id"`
}

// Seat Model
type Seat struct {
	SeatID     int    `json:"seat_id"`
	Number     string `json:"number"`
	QRCodeUrl  string `json:"qr_code_url"`
	MerchantID int    `json:"merchant_id"`
}

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

// OrderHasFood Model
type OrderHasFood struct {
	OrderID int `json:"order_id"`
	FoodID  int `json:"food_id"`
	Amount  int `json:"amount"`
}

// Customer Model
type Customer struct {
	CustomerID int     `json:"customer_id"`
	WechatID   string  `json:"wechat_id"`
	Balance    float32 `json:"balance"`
}
