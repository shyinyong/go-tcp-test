package main

import (
	"github.com/shyinyong/go-tcp-test/server/login"
)

func main() {
	// Login
	loginServer := login.NewLoginServer()
	loginServer.Start("localhost:8081") // Change the port as needed
}
