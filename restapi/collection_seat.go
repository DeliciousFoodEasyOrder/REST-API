package restapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routeSeatCollection(router *mux.Router) {
	base := "/seats"

	router.HandleFunc(base, handlerSecure(handlerListSeats())).
		Methods(http.MethodGet)

	router.HandleFunc(base, handlerSecure(handlerCreateSeat())).
		Methods(http.MethodPost)

	router.HandleFunc(base+"/{ID}", handlerSecure(handlerDeleteSeat())).
		Methods(http.MethodDelete)
}

func handlerListSeats() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

	}

}

func handlerCreateSeat() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

	}

}

func handlerDeleteSeat() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

	}

}
