package models

// Food Model
type Food struct {
	FoodID      int     `json:"food_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	MerchantID  int     `json:"merchant_id"`
}
