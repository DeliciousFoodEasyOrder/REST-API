package restapi

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/DeliciousFoodEasyOrder/REST-API/models"
	"github.com/gorilla/mux"
)

func routeMerchantCollection(router *mux.Router) {
	base := "/merchants"

	// ### Get a merchant [GET /merchants{?email,phone}]
	router.HandleFunc(base, handlerSecure(handlerGetMerchant())).
		Methods(http.MethodGet)

	// ### Create a merchant [POST /merchants]
	router.HandleFunc(base, handlerCreateMerchant()).
		Methods(http.MethodPost)

	// ### Update a merchant partially [PATCH /merchants/{id}]
	router.HandleFunc(base+"/{id}", handlerSecure(handlerPatchMerchant())).
		Methods(http.MethodPatch)
}

func handlerPatchMerchant() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		merchantIDStr := mux.Vars(req)["id"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", merchantIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"修改商户失败",
				NewErr("Bad parameters", "id must be a number"),
			))
			return
		}

		merchantID, _ := strconv.Atoi(merchantIDStr)
		var old *models.Merchant
		if old = models.MerchantDAO.FindByID(merchantID); old == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"修改商户失败",
				NewErr("Bad parameters", "merchant not found"),
			))
			return
		}

		var data models.Merchant
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&data)
		data.MerchantID = merchantID
		merchant, err := models.MerchantDAO.UpdateOne(&data)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"修改商户失败",
				NewErr("Database error", "see server log for more information"),
			))
			panic(err)
		}
		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"修改商户成功",
			merchant,
		))

	}

}

func handlerGetMerchant() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		email := req.FormValue("email")
		phone := req.FormValue("phone")

		var merchant *models.Merchant
		if email != "" {
			merchant = models.MerchantDAO.FindByEmail(email)
		} else if phone != "" {
			merchant = models.MerchantDAO.FindByPhone(phone)
		} else {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取商户信息失败",
				NewErr("Bad parameters", "parameters cannot all be empty"),
			))
			return
		}
		if merchant == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取商户信息失败",
				NewErr("User not found", ""),
			))
			return
		}

		formatter.JSON(w, http.StatusOK, NewResp(
			http.StatusOK,
			"获取商户信息成功",
			merchant,
		))

	}

}

func handlerCreateMerchant() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		var merchant models.Merchant
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&merchant)
		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"商户注册失败",
				NewErr("Bad data", "data format may be incorrect"),
			))
			panic(err)
		}

		emailReg := "^[a-zA-Z0-9_.-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$"
		phoneReg := "^1[0-9]{10}$"
		validEmail, _ := regexp.MatchString(emailReg, merchant.Email)
		validPhone, _ := regexp.MatchString(phoneReg, merchant.Phone)
		if !validEmail || !validPhone {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"商户注册失败",
				NewErr("Bad data", "email or phone invalid"),
			))
			return
		}

		var m *models.Merchant
		m = models.MerchantDAO.FindByEmail(merchant.Email)
		m = models.MerchantDAO.FindByPhone(merchant.Phone)
		if m != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"商户注册失败QAQ",
				NewErr("Bad data", "email or phone already exist"),
			))
			return
		}

		err = models.MerchantDAO.InsertOne(&merchant)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"商户注册失败",
				NewErr("Database error", "see server log for more information"),
			))
			return
		}
		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"商户注册成功",
			merchant,
		))
	}

}
