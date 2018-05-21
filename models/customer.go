package models

// Customer Model
type Customer struct {
	CustomerID int     `json:"customer_id"`
	WechatID   string  `json:"wechat_id"`
	Balance    float32 `json:"balance"`
}

// CustomerDataAccessObject provides access for Model Customer
type CustomerDataAccessObject struct{}

// CustomerDAO is an instance of CustomerDataAccessObject
var CustomerDAO *CustomerDataAccessObject
