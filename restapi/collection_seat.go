package restapi

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/skip2/go-qrcode"

	"github.com/DeliciousFoodEasyOrder/REST-API/models"

	"github.com/DeliciousFoodEasyOrder/REST-API/token"
	"github.com/gorilla/mux"
)

func routeSeatCollection(router *mux.Router) {
	base := "/seats"

	// ### List seats [GET /seats{?merchant_id}]
	router.HandleFunc(base, handlerSecure(handlerListSeats())).
		Methods(http.MethodGet)

	// ### Create a seat [POST /seats]
	router.HandleFunc(base, handlerSecure(handlerCreateSeat())).
		Methods(http.MethodPost)

	// ### Delete a seat [DELETE /seats/{seat_id}]
	router.HandleFunc(base+"/{ID}", handlerSecure(handlerDeleteSeat())).
		Methods(http.MethodDelete)
}

func handlerListSeats() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		claims := token.ParseClaims(getTokenString(req))
		merchantIDStr := req.FormValue("merchant_id")
		merchantID, _ := strconv.Atoi(merchantIDStr)
		if cond, _ := regexp.MatchString("[1-9][0-9]*", merchantIDStr); !cond {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"获取座位信息失败",
				NewErr("Bad parameters", "merchant_id must be a number"),
			))
			return
		}
		if merchantID != int(claims["aud"].(float64)) {
			formatter.JSON(w, http.StatusUnauthorized, NewResp(
				http.StatusUnauthorized,
				"获取座位列表失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

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
		}

		if seat.MerchantID != int(claims["aud"].(float64)) {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建座位失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		if s := models.SeatDAO.FindByMerchantIDAndNumber(seat.MerchantID,
			seat.Number); s != nil {
			formatter.JSON(w, http.StatusBadRequest, NewResp(
				http.StatusBadRequest,
				"创建座位失败",
				NewErr("Bad parameters", "seat already exist"),
			))
			return
		}

		err = models.SeatDAO.InsertOne(&seat)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建座位失败",
				NewErr("Database error", "see server log for more information"),
			))
			panic(err)
		}
		qrCodePath := "static/qrcodes/" + strconv.Itoa(seat.MerchantID) + "/" +
			strconv.Itoa(seat.SeatID) + ".png"
		os.MkdirAll("static/qrcodes/"+strconv.Itoa(seat.MerchantID), os.ModePerm)
		seat.QRCodeURL = "/" + qrCodePath
		seatJSON, _ := json.Marshal(seat)
		err = qrcode.WriteFile(string(seatJSON), qrcode.Medium, 256, qrCodePath)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, NewResp(
				http.StatusInternalServerError,
				"创建座位失败",
				NewErr("QRCode error", "see server log for more information"),
			))
			models.SeatDAO.DeleteBySeatID(seat.SeatID)
			panic(err)
		}
		models.SeatDAO.UpdateOne(&seat)

		formatter.JSON(w, http.StatusCreated, NewResp(
			http.StatusCreated,
			"创建座位成功",
			seat,
		))
	}
}

func handlerDeleteSeat() http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		seatIDStr := mux.Vars(req)["ID"]
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

		if seat.MerchantID != int(claims["aud"].(float64)) {
			formatter.JSON(w, http.StatusUnauthorized, NewResp(
				http.StatusUnauthorized,
				"删除座位失败",
				NewErr("Permission denied", "id mismatch"),
			))
			return
		}

		models.SeatDAO.DeleteBySeatID(seatID)
		formatter.JSON(w, http.StatusNoContent, nil)
	}

}
