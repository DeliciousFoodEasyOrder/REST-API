package main

import "github.com/DeliciousFoodEasyOrder/REST-API/restapi"

func main() {
	port := ":80"

	server := restapi.NewServer()

	server.Run(port)
}
