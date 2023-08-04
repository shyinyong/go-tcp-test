package handler

import (
	"fmt"
	"net"
)

func HandleMsg(conn net.Conn, msg []byte) {
	fmt.Printf("Hello:[%s]\n", string(msg))

	//var gameMsg gen.GameMsg
	//err := proto.Unmarshal(msg, &gameMsg)
	//if err != nil {
	//	// Handle error
	//	return
	//}

	// Handle the message
	// This is where you would implement your game logic
	// For example, you could use gameMsg.GetCommand() to determine what to do
	// And you could use mysql.DB and redis.RDB to interact with your databases
}
