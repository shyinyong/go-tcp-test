package main

import (
	"github.com/shyinyong/go-tcp-test/server/gateway"
	"log"
)

func main() {
	// Gateway
	gateway := gateway.NewGatewayServer()
	// Configure server addresses
	gateway.AddServers()
	//gateway.AddServerAddr(protobuf.MessageType_Login, "login-server:12346")    // Change the address
	// Establish connections to game servers
	err := gateway.ConnectToServers()
	if err != nil {
		log.Fatal("Error connecting to servers:", err)
	}
	gateway.Start("localhost:8080") // Change the port as needed
}
