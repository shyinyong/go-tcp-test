package main

import (
	"github.com/shyinyong/go-tcp-test/handler/login"
)

func main() {
	// Login
	loginServer := login.NewServer()
	loginServer.Start("localhost:8081") // Change the port as needed
}
