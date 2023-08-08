package main

import (
	"github.com/shyinyong/go-tcp-test/handler/gateway"
)

func main() {
	// Create a new instance of the gateway server
	gs := gateway.NewServer()

	// Initialize the serverAddresses map with the addresses of different servers
	//serverAddresses := map[string]string{
	//	"login": "localhost:8081",
	//	//"game":  "game-server:12347",
	//	//"chat":  "chat-server:12348",
	//	// Add more server addresses if needed...
	//}
	//gs.SetServerAddresses(serverAddresses)
	gs.Start("localhost:8080")
}
