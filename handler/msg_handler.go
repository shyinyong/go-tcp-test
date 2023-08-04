package handler

import (
	"fmt"
	"net"
	"sync"
)

type OnlinePlayers struct {
	sync.Mutex
	players map[int32]PlayerInfo
}

type PlayerInfo struct {
	Username string
	ClientIP string
}

func HandleMsg(conn net.Conn, msg []byte) {
	onlinePlayers := &OnlinePlayers{
		players: make(map[int32]PlayerInfo),
	}
	// Add player to online players list
	onlinePlayers.Lock()
	onlinePlayers.players[1] = PlayerInfo{}
	onlinePlayers.Unlock()

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
