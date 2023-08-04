package handler

import (
	"fmt"
	"github.com/shyinyong/go-tcp-test/pb/request"
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

func HandleMsg(conn net.Conn, msg *request.Request) {
	onlinePlayers := &OnlinePlayers{
		players: make(map[int32]PlayerInfo),
	}
	// Add player to online players list
	onlinePlayers.Lock()
	onlinePlayers.players[1] = PlayerInfo{}
	onlinePlayers.Unlock()
	id := msg.Id

	fmt.Printf("Hello:[request id:%s]\n", string(id))

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

	// step1 https://protobuf.dev/getting-started/gotutorial/
	// step2 go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	// step3 protoc -I=D:\projects\go-tcp-test --go_out=D:\projects\go-tcp-test\pb D:\projects\go-tcp-test\proto\*.proto

	// 在读取请求消息的逻辑中，首先读取消息头中的长度字段，然后根据长度字段读取消息体。
	// 在发送响应消息的逻辑中，先构造消息头，然后将消息头和消息体合并后一次性发送。
	// 注意在发送消息时，需要将消息头和消息体合并为一个字节数组进行发送。
}
