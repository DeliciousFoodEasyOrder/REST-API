package main

import (
	"github.com/DeliciousFoodEasyOrder/REST-API/restapi"
)

func main() {
	port := ":8080"

	server := restapi.NewServer()

	// server.Run(port)
	server.Run(port)
}
