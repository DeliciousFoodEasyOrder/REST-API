package restapi

import (
	"net/http"
	"regexp"

	"github.com/DeliciousFoodEasyOrder/REST-API/models"
	"github.com/DeliciousFoodEasyOrder/REST-API/token"
	"github.com/gorilla/mux"
)

func routeOAuthCollection(router *mux.Router) {
	base := "/oauth"

	// ### Password Grant [POST /oauth/token{?grant_type,username,password}]
	router.HandleFunc(base+"/token", handlerPasswordGrant()).
		Methods(http.MethodPost)
}

func handlerPasswordGrant() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		grantType := req.FormValue("password")
		username := req.FormValue("username")
		password := req.FormValue("password")

		if grantType != "password" {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"认证失败",
				NewErr("grant_type error", "grant_type must be password"),
			))
			return
		}

		if username == "" || password == "" {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"认证失败",
				NewErr("Bad parameters", "username or password invalid"),
			))
			return
		}

		emailReg := "^[a-zA-Z0-9_-.]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$"
		phoneReg := "^1[0-9]{10}$"
		isEmail, _ := regexp.MatchString(emailReg, username)
		isPhone, _ := regexp.MatchString(phoneReg, username)
		var merchant *models.Merchant
		if isEmail {
			merchant = models.MerchantDAO.FindByEmail(username)
		} else if isPhone {
			merchant = models.MerchantDAO.FindByPhone(username)
		} else {
			merchant = nil
		}
		if merchant == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"认证失败",
				NewErr("Bad parameters", "username not found"),
			))
			return
		}
		if merchant.Password != password {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"认证失败",
				NewErr("Permission denied", "username or password incorrect"),
			))
			return
		}

		formatter.JSON(w, http.StatusOK, token.NewJWTToken(
			merchant.MerchantID, 259200))
	}

}
