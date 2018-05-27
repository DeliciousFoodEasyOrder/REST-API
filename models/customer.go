package models

// Customer Model
type Customer struct {
	CustomerID int     `xorm:"PK AUTOINCR" json:"customer_id"`
	WechatID   string  `json:"wechat_id"`
	Balance    float32 `json:"balance"`
}

// CustomerDataAccessObject provides access for Model Customer
type CustomerDataAccessObject struct{}

// CustomerDAO is an instance of CustomerDataAccessObject
var CustomerDAO *CustomerDataAccessObject

// InsertOne inserts a customer to database
func (*CustomerDataAccessObject) InsertOne(c *Customer) error {
	_, err := orm.InsertOne(c)
	return err
}

// FindByCustomerID finds a Customer by its ID
func (*CustomerDataAccessObject) FindByCustomerID(customerID int) *Customer {
	var customer Customer
	has, err := orm.Table(customer).ID(customerID).Get(&customer)
	
	if err != nil {
		panic(err)
	}
	if !has {
		return nil
	}
	return &customer
}

// UpdateOne updates an customer
func (*CustomerDataAccessObject) UpdateOne(customer *Customer) (*Customer, error) {
	_, err := orm.Table(customer).ID(customer.CustomerID).Update(customer)
	if err != nil {
		return nil, err
	}
	return CustomerDAO.FindByCustomerID(customer.CustomerID), nil
}