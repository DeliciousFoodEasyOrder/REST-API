package models

// Merchant Model
type Merchant struct {
	MerchantID int    `json:"merchant_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
}

// MerchantDataAccessObject provides access for Model Merchant
type MerchantDataAccessObject struct{}

// MerchantDAO is an instance of MerchantDataAccessObject
var MerchantDAO *MerchantDataAccessObject

// InsertOne inserts a merchant to database
func (*MerchantDataAccessObject) InsertOne(m *Merchant) {
	_, err := orm.InsertOne(m)
	if err != nil {
		panic(err)
	}
}

// FindByEmail finds a merchant by email
func (*MerchantDataAccessObject) FindByEmail(email string) *Merchant {
	var merchant Merchant
	_, err := orm.Table(merchant).Where("Email=?", email).Get(&merchant)
	if err != nil {
		panic(err)
	}
	return &merchant
}

// FindByPhone finds a merchant by phone
func (*MerchantDataAccessObject) FindByPhone(phone string) *Merchant {
	var merchant Merchant
	_, err := orm.Table(merchant).Where("Phone=?", phone).Get(&merchant)
	if err != nil {
		panic(err)
	}
	return &merchant
}
