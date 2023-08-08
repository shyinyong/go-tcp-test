package main

import (
	"github.com/shyinyong/go-tcp-test/handler/gateway"
	"log"
)

func main() {
	// Create a new instance of the gateway server
	gs := gateway.NewServer()

	// Initialize the serverAddresses map with the addresses of different servers
	serverAddresses := map[string]string{
		"login": "localhost:8081",
		//"game":  "game-server:8082",
		//"chat":  "chat-server:8083",
		// Add more server addresses if needed...
	}
	gs.SetServerAddresses(serverAddresses)

	// Start the gateway server and listen on multiple addresses
	addresses := []string{"localhost:8080"}
	err := gs.Start(addresses)
	if err != nil {
		log.Fatal("Error starting gateway server:", err)
	}
}
