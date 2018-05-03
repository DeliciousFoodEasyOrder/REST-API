package models

// Merchant Model
type Merchant struct {
	MerchantID int    `json:"merchant_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
}
