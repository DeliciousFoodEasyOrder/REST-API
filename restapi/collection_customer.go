package restapi

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/DeliciousFoodEasyOrder/REST-API/models"
	"github.com/gorilla/mux"
)

func routeCustomerCollection(router *mux.Router) {
	base := "/customers"

	// ### Get a Customer [GET /customers/{customer_id}]
	router.HandleFunc(base+"/{customer_id}", handlerGetCustomer()).
		Methods(http.MethodGet)

	// ### Create a customer [POST /customers]
	router.HandleFunc(base, handlerCreateCustomer()).
		Methods(http.MethodPost)

	// ### Update a customer [PUT /customers/{customer_id}]
	router.HandleFunc(base+"/{customer_id}", handlerPutCustomer()).
		Methods(http.MethodPatch)

}

func handlerGetCustomer() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		customerIDStr := mux.Vars(req)["customer_id"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", customerIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取用户信息失败",
				NewErr("Bad parameters", "id must be a number"),
			))
			return
		}

		customerID, _ := strconv.Atoi(customerIDStr)
		customer := models.CustomerDAO.FindByCustomerID(customerID)
		if customer == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取用户信息失败",
				NewErr("Bad parameters", "customer not found"),
			))
			return
		}

		formatter.JSON(w, http.StatusOK, NewResp(
			http.StatusOK,
			"获取用户信息成功",
			customer,
		))
	}
}

func handlerCreateCustomer() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		//claims := token.ParseClaims(getTokenString(req))
		decoder := json.NewDecoder(req.Body)
		var customer models.Customer
		err := decoder.Decode(&customer)
		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建用户失败",
				NewErr("Bad parameters", "please check your request format"),
			))
			panic(err)
		}

		err = models.CustomerDAO.InsertOne(&customer)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建用户失败",
				NewErr("Database error", "see server log for more information"),
			))
			panic(err)
		}
		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"创建用户成功",
			customer,
		))

	}
}

func handlerPutCustomer() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		customerIDStr := mux.Vars(req)["customer_id"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", customerIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取用户信息失败",
				NewErr("Bad parameters", "id must be a number"),
			))
			return
		}

		customerID, _ := strconv.Atoi(customerIDStr)
		var old *models.Customer

		if old = models.CustomerDAO.FindByCustomerID(customerID); old == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"修改用户信息失败",
				NewErr("Bad parameters", "customer not found"),
			))
			return
		}

		var data models.Customer
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&data)

		customer, err := models.CustomerDAO.UpdateOne(&data)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"修改用户信息失败",
				NewErr("Database error", "see server log for more information"),
			))
			panic(err)
		}
		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"修改用户信息成功",
			customer,
		))
	}

}
