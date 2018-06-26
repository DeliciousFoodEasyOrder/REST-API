package models

// Food Model
type Food struct {
	FoodID      int     `xorm:"PK AUTOINCR" json:"food_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	MerchantID  int     `json:"merchant_id"`
	IconURL     string  `json:"icon_url"`
}

// FoodWithAmount specifies the amount in an specific order
type FoodWithAmount struct {
	Food
	Amount int `json:"amount"`
}

// FoodDataAccessObject provides access for Model Food
type FoodDataAccessObject struct{}

// FoodDAO is an instance of FoodDataAccessObject
var FoodDAO *FoodDataAccessObject

// InsertOne inserts a food to database
func (*FoodDataAccessObject) InsertOne(food *Food) error {
	_, err := orm.InsertOne(food)
	return err
}

// FindByID finds a food by id
func (*FoodDataAccessObject) FindByID(foodID int) *Food {
	var food Food
	has, err := orm.Table(food).ID(foodID).Get(&food)
	if err != nil {
		panic(err)
	}
	if !has {
		return nil
	}
	return &food
}

// FindByMerchantID finds foods by a merchant ID
func (*FoodDataAccessObject) FindByMerchantID(merchantID int) []Food {
	foods := make([]Food, 0)
	err := orm.Table("Food").Where("MerchantID=?", merchantID).Find(&foods)
	if err != nil {
		panic(err)
	}
	return foods
}

// DeleteByFoodID deletes a food by id
func (*FoodDataAccessObject) DeleteByFoodID(foodID int) {
	var food Food
	_, err := orm.Table(food).ID(foodID).Delete(&food)
	if err != nil {
		panic(err)
	}
}

// UpdateOne updates a food by id of parameter
func (*FoodDataAccessObject) UpdateOne(food *Food) error {
	_, err := orm.Table(food).ID(food.FoodID).Update(food)
	return err
}
