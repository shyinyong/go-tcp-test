package handler

import (
	"fmt"
	"github.com/shyinyong/go-tcp-test/pb/request"
	"google.golang.org/protobuf/proto"
	"net"
)

func HandleMsg(conn net.Conn, msg *request.Request) {
	// Handle different message IDs
	switch msg.MsgId {
	case 1: // 用户信息请求
		response, err = GetUserInfo(&msg)
		if err != nil {
			fmt.Println("Error handling user info request:", err)
			return
		}
		// TODO: Handle other message IDs

	default:
		fmt.Println("Unknown message ID:", request.MsgID)
		return
	}

	// Serialize response message
	responseData, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		return
	}

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
