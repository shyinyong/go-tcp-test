package main

import "github.com/shyinyong/go-tcp-test/handler/chat"

func main() {
	chatServer := chat.NewServer()
	chatServer.Start("localhost:8083")
}
