package models

// Merchant Model
type Merchant struct {
	MerchantID int    `xorm:"PK AUTOINCR" json:"merchant_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	On         int    `json:"on"`
}

// MerchantDataAccessObject provides access for Model Merchant
type MerchantDataAccessObject struct{}

// MerchantDAO is an instance of MerchantDataAccessObject
var MerchantDAO *MerchantDataAccessObject

// InsertOne inserts a merchant to database
func (*MerchantDataAccessObject) InsertOne(m *Merchant) error {
	_, err := orm.InsertOne(m)
	return err
}

// UpdateOne updates a merchant
func (*MerchantDataAccessObject) UpdateOne(m *Merchant) (*Merchant, error) {
	_, err := orm.Id(m.MerchantID).Update(m)
	if err != nil {
		return nil, err
	}
	return MerchantDAO.FindByID(m.MerchantID), nil
}

// FindByID finds a merchant by MerchantID
func (*MerchantDataAccessObject) FindByID(merchantID int) *Merchant {
	var merchant Merchant
	has, err := orm.Table(merchant).ID(merchantID).Get(&merchant)
	if err != nil {
		panic(err)
	}
	if !has {
		return nil
	}
	return &merchant
}

// FindByEmail finds a merchant by email
func (*MerchantDataAccessObject) FindByEmail(email string) *Merchant {
	var merchant Merchant
	has, err := orm.Table(merchant).Where("Email=?", email).Get(&merchant)
	if err != nil {
		panic(err)
	}
	if !has {
		return nil
	}
	return &merchant
}

// FindByPhone finds a merchant by phone
func (*MerchantDataAccessObject) FindByPhone(phone string) *Merchant {
	var merchant Merchant
	has, err := orm.Table(merchant).Where("Phone=?", phone).Get(&merchant)
	if err != nil {
		panic(err)
	}
	if !has {
		return nil
	}
	return &merchant
}
