package main

import (
	"github.com/shyinyong/go-tcp-test/server/gateway"
	"log"
)

func main() {
	// Gateway
	gatewayServer := gateway.NewGatewayServer()
	// Configure server addresses
	gatewayServer.AddServers()
	//gatewayServer.AddServerAddr(protobuf.MessageType_Login, "login-server:12346")    // Change the address
	// Establish connections to game servers
	err := gatewayServer.ConnectToServers()
	if err != nil {
		log.Fatal("Error connecting to servers:", err)
	}
	gatewayServer.Start("localhost:8080") // Change the port as needed
}
