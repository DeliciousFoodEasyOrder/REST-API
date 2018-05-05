package restapi

import (
	"encoding/json"
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
		claims := token.ParseClaims(getTokenString(req))
		merchantIDStr := req.FormValue("merchant_id")
		var []seatList *models.Seat
		if cond, _ := regexp.MatchString("[1-9][0-9]*", merchantIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取座位信息失败",
				NewErr("Bad parameters", "merchant_id must be a number"),
			))
			return
		}
		if merchantIDStr != claims["aud"] {
			formatter.JSON(w, http.StatusUnauthorized, NewResp(
				http.StatusUnauthorized,
				"获取座位列表失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		var merchantID int
		merchantID, _ = strconv.Atoi(merchantIDStr)
		seatList := models.SeatDAO.FindByMerchantID(merchantID)
		formatter.JSON(w, http.StatusOK, NewResp(
			http.StatusOK,
			"获取座位列表成功",
			seatList,
		))
	}
}

func handlerCreateSeat() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		claims := token.ParseClaims(getTokenString(req))
		decoder := json.NewDecoder(req.Body)
		var seat models.Seat
		err := decoder.Decode(&seat)
		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建座位失败",
				NewErr("Bad parameters", "please check your request format"),
			))
			panic(err)
			return
		}

		if seat.SeatID != claims["aud"] {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建座位失败",
				NewErr("Permission denied", "id mismatch"),
			))
			panic(err)
			return
		}

		err = models.SeatDAO.InsertOne(&order)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建座位失败",
				NewErr("Database error", "see server log for more information"),
			))
			panic(err)
			return
		}
		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"创建座位成功",
			seat,
		))
    }
}

func handlerDeleteSeat() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		//req.ParseForm()
		seatIDStr := mux.Vars(req)["ID"]
		//seatID, _ :=  strconv.Atoi(mux.Vars(req)["ID"])
		if cond, _ := regexp.MatchString("[1-9][0-9]*", seatIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"删除座位失败",
				NewErr("Bad parameters", "ID must be a number"),
			))
			return
		}

		seatID, _ := strconv.Atoi(seatIDStr)
		claims := token.ParseClaims(getTokenString(req))
		seat := models.SeatDAO.FindByID(seatID)

		if seat == nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"删除座位失败",
				NewErr("Bad parameters", "seat not found"),
			))
			return
		}

		if seat.MerchantID != claims["aud"] {
			formatter.JSON(w, http.StatusUnauthorized, NewResp(
				http.StatusUnauthorized,
				"删除座位失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		models.SeatDao.DeleteBySeatID(seatID)
		formatter.JSON(w, http.StatusNoContent, nil)
	}

}
