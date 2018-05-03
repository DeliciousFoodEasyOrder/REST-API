package models

// Customer Model
type Customer struct {
	CustomerID int     `json:"customer_id"`
	WechatID   string  `json:"wechat_id"`
	Balance    float32 `json:"balance"`
}
