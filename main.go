package main

import (
	"net/http"

	"github.com/DeliciousFoodEasyOrder/REST-API/restapi"
)

func main() {
	port := ":443"

	server := restapi.NewServer()

	// server.Run(port)
	http.ListenAndServeTLS(
		port,
		"/etc/letsencrypt/live/easyorder.cf/fullchain.pem",
		"/etc/letsencrypt/live/easyorder.cf/privkey.pem",
		server,
	)
}
