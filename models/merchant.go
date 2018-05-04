package models

// Merchant Model
type Merchant struct {
	MerchantID int    `json:"merchant_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
}

// Insert inserts self to database
func (m *Merchant) Insert() {
	_, err := orm.InsertOne(m)
	if err != nil {
		panic(err)
	}
}

// FindByEmail finds a merchant by email
func FindByEmail(email string) *Merchant {
	var merchant Merchant
	_, err := orm.Table(merchant).Where("Email=?", email).Get(&merchant)
	if err != nil {
		panic(err)
	}
	return &merchant
}

// FindByPhone finds a merchant by phone
func FindByPhone(phone string) *Merchant {
	var merchant Merchant
	_, err := orm.Table(merchant).Where("Phone=?", phone).Get(&merchant)
	if err != nil {
		panic(err)
	}
	return &merchant
}
