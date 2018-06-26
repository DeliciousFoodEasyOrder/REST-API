package restapi

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/DeliciousFoodEasyOrder/REST-API/models"

	"github.com/DeliciousFoodEasyOrder/REST-API/token"
	"github.com/gorilla/mux"

	"path/filepath"
	"os"
	"io/ioutil"
)

const maxUploadSize = 4 * 1024 * 1024
const foodsPath = "static/foods/"

func routeFoodCollection(router *mux.Router) {
	base := "/foods"

	// ### List foods [GET /foods{?merchant_id}]
	router.HandleFunc(base, handlerListFoods()).
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

	// ### Create a icon of a food [POST /foods/{food_id}/icon]
	router.HandleFunc(base+"/{food_id}/icon", handlerSecure(handlerCreateIconOfFood())).
	    Methods(http.MethodPost)
}

func handlerCreateIconOfFood() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		foodIDStr := mux.Vars(req)["food_id"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", foodIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取对应菜品失败",
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
				"获取对应菜品失败",
				NewErr("Bad parameters", "food not found"),
			))
			return
		}

		if food.MerchantID != int(claims["aud"].(float64)) {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取菜品失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)
		if err := req.ParseMultipartForm(maxUploadSize); err != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建图片失败",
				NewErr("FILE_TOO_BIG", "PLEASE MINIFY IT"),
			))
			panic(err)
		}

		//fileType := req.PostFormValue("type") 这个是有type框的时候这么写
		//要这么写吗
		file, _, err := req.FormFile("uploadFile")
		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建图片失败",
				NewErr("INVALID_FILE", "PLEASE MODIFY IT"),
			))
			panic(err)
		}

		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建图片失败",
				NewErr("INVALID_FILE", "PLEASE MODIFY IT"),
			))
			panic(err)
		}

		fileType := http.DetectContentType(fileBytes)
		if fileType != "image/jpeg" && fileType != "image/jpg" &&
		    fileType != "image/gif" && fileType != "image/png" {
				formatter.JSON(w, http.StatusBadRequest, NewResp(
					http.StatusBadRequest,
					"创建图片失败",
					NewErr("INVALID_FILE_TYPE", "PLEASE MODIFY THE TYPE"),
				))
				panic(err)
		}

		// 向static/foods/下写文件
		fileName := foodIDStr
		newFoodsPath := filepath.Join(foodsPath, fileName)
		newFile, err := os.Create(newFoodsPath)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建图片失败",
				NewErr("CANNOT_WRITE_FILE_TO_FOODS", "PLEASE MODIFY IT"),
			))
			panic(err)
		}

		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建图片失败",
				NewErr("CANNOT_WRITE_FILE_TO_FOODS", "PLEASE MODIFY IT"),
			))
			panic(err)
		}

		// 将图片的url插入数据库(food表)
		food.IconURL = filepath.Join("/", newFoodsPath)
		err = models.FoodDAO.UpdateOne(food)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建图片失败",
				NewErr("Database error", "see server log for more details"),
			))
			panic(err)
		}

		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusOK,
			"创建图片成功",
			food,
		))
	}
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

		if food.MerchantID != int(claims["aud"].(float64)) {
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
			return
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

		if food.MerchantID != int(claims["aud"].(float64)) {
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
		merchantID, _ := strconv.Atoi(merchantIDStr)

		if cond, _ := regexp.MatchString("[1-9][0-9]*", merchantIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取菜品列表失败",
				NewErr("Bad parameters", "merchant_id must be a number"),
			))
			return
		}

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
