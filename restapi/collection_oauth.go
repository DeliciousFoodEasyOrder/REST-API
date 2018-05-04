package restapi

import (
	"net/http"

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

		// TODO: Find user in database

		if !(username == "a" && password == "b") {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"认证失败",
				NewErr("Permission denied", "username or password incorrect"),
			))
			return
		}

		formatter.JSON(w, http.StatusOK, token.NewJWTToken(1, 259200))
	}

}
