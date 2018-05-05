package restapi

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/DeliciousFoodEasyOrder/REST-API/models"

	"github.com/DeliciousFoodEasyOrder/REST-API/token"
	"github.com/gorilla/mux"
)

func routeOrderCollection(router *mux.Router) {
	base := "/orders"

	// ### List orders [GET /orders{?merchant_id,status}]
	router.HandleFunc(base, handlerSecure(handlerListOrders())).
		Methods(http.MethodGet)

	// ### Get an order [GET /orders/{id}]
	router.HandleFunc(base+"/{ID}", handlerSecure(handlerGetOrder())).
		Methods(http.MethodGet)

	// ### Create an order [POST /orders]
	router.HandleFunc(base, handlerSecure(handlerCreateOrder())).
		Methods(http.MethodPost)

	// ### Update an order partially [PATCH /orders/{id}]
	router.HandleFunc(base+"/{ID}", handlerSecure(handlerPatchOrder())).
		Methods(http.MethodPatch)

}

func handlerListOrders() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		claims := token.ParseClaims(getTokenString(req))
		merchantIDStr := req.FormValue("merchant_id")
		statusStr := req.FormValue("status")

		if cond, _ := regexp.MatchString("[1-9][0-9]*", merchantIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取订单列表失败",
				NewErr("Bad parameters", "merchant_id must be a number"),
			))
			return
		}

		if cond, _ := regexp.MatchString("0|1", statusStr); !cond &&
			statusStr != "" {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取订单列表失败",
				NewErr("Bad parameters", "status must be 0, 1 or empty"),
			))
			return
		}

		if merchantIDStr != claims["aud"] {
			formatter.JSON(w, http.StatusUnauthorized, NewResp(
				http.StatusUnauthorized,
				"获取订单列表失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		var merchantID, status int
		merchantID, _ = strconv.Atoi(merchantIDStr)
		if statusStr == "" {
			status = -1
		} else {
			status, _ = strconv.Atoi(statusStr)
		}
		orders := models.OrderDAO.FindByMerchantIDAndStatus(merchantID, status)
		formatter.JSON(w, http.StatusOK, NewResp(
			http.StatusOK,
			"获取订单列表成功",
			orders,
		))
	}

}

func handlerGetOrder() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		orderIDStr := mux.Vars(req)["ID"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", orderIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取订单失败",
				NewErr("Bad parameters", "id must be a number"),
			))
			return
		}

		orderID, _ := strconv.Atoi(orderIDStr)
		order := models.OrderDAO.FindByOrderID(orderID)
		formatter.JSON(w, http.StatusOK, NewResp(
			http.StatusOK,
			"获取订单成功",
			order,
		))
	}

}

func handlerCreateOrder() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		claims := token.ParseClaims(getTokenString(req))
		decoder := json.NewDecoder(req.Body)
		var order models.OrderWithFoods
		err := decoder.Decode(&order)
		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建订单失败",
				NewErr("Bad parameters", "please check your request format"),
			))
			panic(err)
			return
		}

		if order.MerchantID != claims["aud"] {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建订单失败",
				NewErr("Permission denied", "id mismatch"),
			))
			panic(err)
			return
		}

		err = models.OrderDAO.InsertOne(&order)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建订单失败",
				NewErr("Database error", "see server log for more information"),
			))
			panic(err)
			return
		}

		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"创建订单成功",
			order,
		))
	}

}

func handlerPatchOrder() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		orderIDStr := mux.Vars(req)["ID"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", orderIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取订单失败",
				NewErr("Bad parameters", "id must be a number"),
			))
			return
		}

		orderID, _ := strconv.Atoi(orderIDStr)
		var old *models.OrderWithFoods

		if old = models.OrderDAO.FindByOrderID(orderID); old == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"修改订单失败",
				NewErr("Bad parameters", "order not found"),
			))
			return
		}

		var data map[string]string
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&data)
		for key, val := range data {
			switch key {
			case "status":
				old.Status, _ = strconv.Atoi(val)
				break
			case "seat_id":
				old.SeatID, _ = strconv.Atoi(val)
				break
			case "customer_id":
				old.CustomerID, _ = strconv.Atoi(val)
				break
			case "order_time":
				old.OrderTime, _ = time.Parse(time.RFC3339, val)
				break
			case "complete_time":
				old.CompleteTime, _ = time.Parse(time.RFC3339, val)
				break
			default:
				formatter.JSON(w, http.StatusBadRequest, NewResp(
					http.StatusBadRequest,
					"修改订单失败",
					NewErr("Bad parameters", `Request contains invalid fields 
or fields cannot be modified`),
				))
				return
			}
		}

		orderWithFoods, err := models.OrderDAO.UpdateOne(&old.Order)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"修改订单失败",
				NewErr("Database error", "see server log for more information"),
			))
			panic(err)
			return
		}
		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"修改订单成功",
			orderWithFoods,
		))
	}

}
