package restapi

import (
	"net/http"

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

	// ### Update a customer partially [PATCH /customers/{customer_id}]
	router.HandleFunc(base+"/{customer_id}", handlerPatchCustomer()).
		Methods(http.MethodPatch)

}

func handlerGetCustomer() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

	}

}

func handlerCreateCustomer() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

	}

}

func handlerPatchCustomer() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

	}

}
