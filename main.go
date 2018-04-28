package main

import "github.com/DeliciousFoodEasyOrder/Restaurant-MS/restapi"

func main() {
	port := ":8080"

	server := restapi.NewServer()

	server.Run(port)
}
