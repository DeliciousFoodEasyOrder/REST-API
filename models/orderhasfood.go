package models

// OrderHasFood Model
type OrderHasFood struct {
	OrderID int `json:"order_id"`
	FoodID  int `json:"food_id"`
	Amount  int `json:"amount"`
}
