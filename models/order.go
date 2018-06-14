package models

import (
	"time"

	"github.com/go-xorm/xorm"
)

// Order Model
type Order struct {
	OrderID      int       `xorm:"PK AUTOINCR" json:"order_id"`
	Status       int       `json:"status"`
	SeatID       int       `json:"seat_id"`
	CustomerID   int       `json:"customer_id"`
	MerchantID   int       `json:"merchant_id"`
	OrderTime    time.Time `json:"order_time"`
	CompleteTime time.Time `json:"complete_time"`
}

// OrderHasFood Model
type OrderHasFood struct {
	ID      int `json:"id PK AUTOINCR"`
	OrderID int `json:"order_id"`
	FoodID  int `json:"food_id"`
	Amount  int `json:"amount"`
}

// OrderWithFoods is the display model of an order
type OrderWithFoods struct {
	Order
	Foods []FoodWithAmount `json:"foods"`
}

// OrderFull joined Order with OrderHasFood and Food
type OrderFull struct {
	Order        `xorm:"extends"`
	OrderHasFood `xorm:"extends"`
	Food         `xorm:"extends"`
}

// OrderDataAccessObject provides database access for Order
type OrderDataAccessObject struct{}

// OrderDAO is an instance of OrderDataAccessObject
var OrderDAO *OrderDataAccessObject

// FindByMerchantIDAndStatus finds orders with merchantID and status specified
// Parameters
// - merchantID : id of a merchant
// - status : status of an order, -1 if no status specified
func (*OrderDataAccessObject) FindByMerchantIDAndStatus(
	merchantID, status int) []OrderWithFoods {

	fullOrders := make([]OrderFull, 0)
	session := OrderDAO.joinFullOrder().Where("Order.MerchantID=?", merchantID)
	if status != -1 {
		session.Where("Status=?", status)
	}
	err := session.Asc("Order.OrderID").Find(&fullOrders)
	if err != nil {
		panic(err)
	}

	orders := make([]OrderWithFoods, 0)
	currentID := 0
	i := -1
	for _, fullOrder := range fullOrders {
		if fullOrder.Order.OrderID != currentID {
			orders = append(orders, OrderWithFoods{
				Order: fullOrder.Order,
				Foods: make([]FoodWithAmount, 0),
			})
			i++
			currentID = fullOrder.Order.OrderID
		}
		orders[i].Foods = append(orders[i].Foods, FoodWithAmount{
			Food:   fullOrder.Food,
			Amount: fullOrder.OrderHasFood.Amount,
		})
	}

	return orders
}

// FindByCustomerIDAndStatus finds orders with customerID and status specified
// Parameters
// - customerID : id of a customer
// - status : status of an order, -1 if no status specified
func (*OrderDataAccessObject) FindByCustomerIDAndStatus(
	customerID, status int) []OrderWithFoods {

	fullOrders := make([]OrderFull, 0)
	session := OrderDAO.joinFullOrder().Where("Order.CustomerID=?", customerID)
	if status != -1 {
		session.Where("Status=?", status)
	}
	err := session.Asc("Order.OrderID").Find(&fullOrders)
	if err != nil {
		panic(err)
	}

	orders := make([]OrderWithFoods, 0)
	currentID := 0
	i := -1
	for _, fullOrder := range fullOrders {
		if fullOrder.Order.OrderID != currentID {
			orders = append(orders, OrderWithFoods{
				Order: fullOrder.Order,
				Foods: make([]FoodWithAmount, 0),
			})
			i++
			currentID = fullOrder.Order.OrderID
		}
		orders[i].Foods = append(orders[i].Foods, FoodWithAmount{
			Food:   fullOrder.Food,
			Amount: fullOrder.OrderHasFood.Amount,
		})
	}

	return orders
}

// FindByOrderID finds an Order by its ID
func (*OrderDataAccessObject) FindByOrderID(orderID int) *OrderWithFoods {
	fullOrders := make([]OrderFull, 0)
	err := OrderDAO.joinFullOrder().Where("Order.OrderID=?", orderID).
		Find(&fullOrders)
	if err != nil {
		panic(err)
	}
	if len(fullOrders) == 0 {
		return nil
	}

	order := &OrderWithFoods{
		Order: fullOrders[0].Order,
		Foods: make([]FoodWithAmount, 0),
	}
	for _, fullOrder := range fullOrders {
		order.Foods = append(order.Foods, FoodWithAmount{
			Food:   fullOrder.Food,
			Amount: fullOrder.OrderHasFood.Amount,
		})
	}

	return order
}

// InsertOne inserts an order into database
func (*OrderDataAccessObject) InsertOne(order *OrderWithFoods) error {
	_, err := orm.InsertOne(&order.Order)
	if err != nil {
		return err
	}
	for _, food := range order.Foods {
		orderHasFood := OrderHasFood{
			OrderID: order.Order.OrderID,
			FoodID:  food.FoodID,
			Amount:  food.Amount,
		}
		_, err := orm.InsertOne(&orderHasFood)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateOne updates an order
func (*OrderDataAccessObject) UpdateOne(order *Order) (*OrderWithFoods, error) {
	_, err := orm.Id(order.OrderID).Update(order)
	if err != nil {
		return nil, err
	}
	return OrderDAO.FindByOrderID(order.OrderID), nil
}

func (*OrderDataAccessObject) joinFullOrder() *xorm.Session {
	return orm.Table("Order").
		Join("LEFT OUTER", "OrderHasFood", "Order.OrderID=OrderHasFood.OrderID").
		Join("LEFT OUTER", "Food", "OrderHasFood.FoodID=Food.FoodID")
}
