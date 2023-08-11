package main

import (
	"github.com/shyinyong/go-tcp-test/config"
	"github.com/shyinyong/go-tcp-test/db/mysql"
	"github.com/shyinyong/go-tcp-test/handler/gateway"
	"log"
)

func main() {
	// Config env initialize
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
		return
	}
	// Initialize database connections
	store := mysql.NewDB(&cfg)

	// Create a new instance of the gateway server
	gs := gateway.NewServer(cfg, store)

	// 创建消息分派器
	dispatcher := gateway.NewDispatcher()
	// 注册处理函数
	dispatcher.RegisterHandler(1, gateway.LoginHandler)
	dispatcher.RegisterHandler(2, gateway.ExitGameHandler)

	// 假设收到消息
	msgType := uint16(1)
	msgID := uint32(1)
	body := []byte("login message body")
	// 调用消息分派器
	dispatcher.Dispatch(msgType, msgID, body)

	// Initialize the serverAddresses map with the addresses of different servers
	serverAddresses := map[string]string{
		"login": "localhost:8081",
		//"game":  "game-server:8082",
		//"chat": "chat-server:8083",
		// Add more server addresses if needed...
	}
	gs.SetServerAddresses(serverAddresses)

	// Start the gateway server and listen on multiple addresses
	addresses := []string{"localhost:8080"}
	err = gs.Start(addresses)
	if err != nil {
		log.Fatal("Error starting gateway server:", err)
	}
}
