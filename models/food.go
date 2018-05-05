package models

// Food Model
type Food struct {
	FoodID      int     `xorm:"AUTOINCR" json:"food_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	MerchantID  int     `json:"merchant_id"`
}

type FoodDataAccessObject struct{}

var FoodDAO *FoodDataAccessObject

// InsertOne inserts a food to database
func (*FoodDataAccessObject) InsertOne(s *Seat) {
	_, err := orm.InsertOne(s)
	if err != nil {
		panic(err)
	}
}

// FindByName finds a food by id
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

func (*FoodDataAccessObject) DeleteByFoodID(foodID int) {
	var food Food
	_, err := orm.Table(food).ID(foodID).Delete(&food)
	if err != nil {
		panic(err)
	}
}

func (*FoodDataAccessObject) UpdateOne(food *Food) {
	_, err := orm.Table(food).ID(food.FoodID).Update(food)
	if err != nil {
		panic(err)
	}
}
