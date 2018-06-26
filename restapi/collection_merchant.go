package restapi

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/DeliciousFoodEasyOrder/REST-API/models"
	"github.com/gorilla/mux"

	"os"
	"path/filepath"
	"io/ioutil"
)

const merchantsPath = "static/merchants/"

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

	// ### Create a icon of a merchant [POST /merchants/{merchant_id}/icon]
	router.HandleFunc(base+"/{merchant_id}/icon", handlerSecure(handlerCreateIconOfMerchant())).
		Methods(http.MethodPost)
		
	// ### Get a merchant by id [GET /mechants/{merchant_id}]
	router.HandleFunc(base+"/{merchant_id}", handlerGetMerchantByID()).
		Methods(http.MethodGet)
}

func handlerGetMerchantByID() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		merchantIDStr := mux.Vars(req)["merchant_id"]

		if cond, _ := regexp.MatchString("[1-9][0-9]*", merchantIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取商家失败",
				NewErr("Bad parameters", "merchant_id must be a number"),
			))
			return
		}

		merchantID, _ := strconv.Atoi(merchantIDStr)
		merchant := models.MerchantDAO.FindByID(merchantID)

		if merchant == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取商家失败",
				NewErr("Bad parameters", "merchant not found"),
			))
			return
		}

		formatter.JSON(w, http.StatusOK, NewResp(
			http.StatusOK,
			"获取商家成功",
			merchant,
		))
	}
}

func handlerCreateIconOfMerchant() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		merchantIDStr := mux.Vars(req)["merchant_id"]


		if cond, _ := regexp.MatchString("[1-9][0-9]*", merchantIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取对应商家失败",
				NewErr("Bad parameters", "merchant_id must be a number"),
			))
			return
		}

		merchantID, _ := strconv.Atoi(merchantIDStr)
		merchant := models.MerchantDAO.FindByID(merchantID)
		if merchant == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取对应商家失败",
				NewErr("Bad parameters", "merchant not found"),
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
		//要这么写
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

		// 向static/merchants/下写文件
		fileName := merchantIDStr
		newMerchantsPath := filepath.Join(merchantsPath, fileName)
		newFile, err := os.Create(newMerchantsPath)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建图片失败",
				NewErr("CANNOT_WRITE_FILE_TO_MERCHANTS", "PLEASE MODIFY IT"),
			))
			panic(err)
		}

		//defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil /*|| newFile.Close() != nil*/ {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建图片失败",
				NewErr("CANNOT_WRITE_FILE_TO_MERCHANTS", "PLEASE MODIFY IT"),
			))
			panic(err)
		}

		if err := newFile.Close(); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建图片失败",
				NewErr("CANNOT_WRITE_FILE_TO_FOODS", "PLEASE MODIFY IT"),
			))
			panic(err)
		}

		// 将图片的url插入数据库(merchant表)
		merchant.IconURL = filepath.Join("/", newMerchantsPath)
		newMerchant, err := models.MerchantDAO.UpdateOne(merchant)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建图片失败",
				NewErr("Database error", "see server log for more details"),
			))
			panic(err)
		}

		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"创建图片成功",
			newMerchant,
		))
	}
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
