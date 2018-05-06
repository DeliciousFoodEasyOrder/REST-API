package restapi

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/DeliciousFoodEasyOrder/REST-API/models"

	"github.com/DeliciousFoodEasyOrder/REST-API/token"
	"github.com/gorilla/mux"
)

func routeFoodCollection(router *mux.Router) {
	base := "/foods"

	// ### List foods [GET /foods{?merchant_id}]
	router.HandleFunc(base, handlerSecure(handlerListFoods())).
		Methods(http.MethodGet)

	// ### Get a food [GET /foods/{food_id}]
	router.HandleFunc(base+"/{food_id}", handlerSecure(handlerGetFood())).
		Methods(http.MethodGet)

	// ### Create a food [POST /foods]
	router.HandleFunc(base, handlerSecure(handlerCreateFood())).
		Methods(http.MethodPost)

	// ### Delete a food [DELETE /foods/{food_id}]
	router.HandleFunc(base+"/{food_id}", handlerSecure(handlerDeleteFood())).
		Methods(http.MethodDelete)
}

func handlerDeleteFood() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		foodIDStr := mux.Vars(req)["food_id"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", foodIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"删除菜品失败",
				NewErr("Bad parameters", "food_id must be a number"),
			))
			return
		}

		foodID, _ := strconv.Atoi(foodIDStr)
		claims := token.ParseClaims(getTokenString(req))
		food := models.FoodDAO.FindByID(foodID)

		if food == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"删除菜品失败",
				NewErr("Bad parameters", "food not found"),
			))
			return
		}

		if food.MerchantID != claims["aud"] {
			formatter.JSON(w, http.StatusUnauthorized, NewResp(
				http.StatusUnauthorized,
				"删除菜品失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		models.FoodDAO.DeleteByFoodID(foodID)

		formatter.JSON(w, http.StatusNoContent, nil)
	}

}

func handlerCreateFood() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		var food models.Food
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&food)
		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建菜品失败",
				NewErr("Bad parameters", "please check your data format"),
			))
			panic(err)
		}

		if models.MerchantDAO.FindByID(food.MerchantID) == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建菜品失败",
				NewErr("Bad parameters", "merchant not found"),
			))
		}

		err = models.FoodDAO.InsertOne(&food)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建菜品失败",
				NewErr("Database error", "see server log for more details"),
			))
			panic(err)
		}

		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusOK,
			"创建菜品成功",
			food,
		))
	}

}

func handlerGetFood() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		foodIDStr := mux.Vars(req)["food_id"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", foodIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取菜品失败",
				NewErr("Bad parameters", "food_id must be a number"),
			))
			return
		}

		foodID, _ := strconv.Atoi(foodIDStr)
		claims := token.ParseClaims(getTokenString(req))
		food := models.FoodDAO.FindByID(foodID)

		if food == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取菜品失败",
				NewErr("Bad parameters", "food not found"),
			))
			return
		}

		if food.MerchantID != claims["aud"] {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取菜品失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		formatter.JSON(w, http.StatusOK, NewResp(
			http.StatusOK,
			"获取菜品成功",
			food,
		))
	}

}

func handlerListFoods() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		merchantIDStr := req.FormValue("merchant_id")
		claims := token.ParseClaims(getTokenString(req))

		if cond, _ := regexp.MatchString("[1-9][0-9]*", merchantIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取菜品列表失败",
				NewErr("Bad parameters", "merchant_id must be a number"),
			))
			return
		}

		if merchantIDStr != claims["aud"] {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取菜品列表失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		merchantID, _ := strconv.Atoi(merchantIDStr)

		if models.MerchantDAO.FindByID(merchantID) == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取菜品列表失败",
				NewErr("Bad parameters", "merchant not found"),
			))
			return
		}

		foods := models.FoodDAO.FindByMerchantID(merchantID)

		formatter.JSON(w, http.StatusOK, NewResp(
			http.StatusOK,
			"获取菜品列表成功",
			foods,
		))
	}

}
