package restapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routeJWTCollection(router *mux.Router) {
	base := "/auth"

	// ### Get a JWT [GET /auth{?username,password,type}]
	router.HandleFunc(base, func(w http.ResponseWriter, req *http.Request) {
		// TODO
	}).Methods(http.MethodGet)
}
