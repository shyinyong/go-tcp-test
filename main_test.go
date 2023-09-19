package main

import (
	"fmt"
	"github.com/shyinyong/go-tcp-test/consts"
	"github.com/shyinyong/go-tcp-test/msg_packet"
	"github.com/shyinyong/go-tcp-test/pb/cs"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"net"
	"testing"
)

func TestRunServer(t *testing.T) {
	// 启动服务器
	gameServer, err := NewGameServer()
	if err != nil {
		t.Fatalf("Error creating GameServer: %v", err)
	}
	// defer gameServer.listener.Close()
	go gameServer.Run()
}

func TestLoginAndPing(t *testing.T) {
	// 启动服务器
	gameServer, err := NewGameServer()
	assert.NoError(t, err)
	gameServer.RegisterHandler(consts.CMMessageID_CMPing, 1)
	gameServer.RegisterHandler(consts.CMMessageID_CMLogin, 1)
	go gameServer.Run()

	// 模拟与服务器建立连接
	conn, err := net.Dial("tcp", "localhost:8080")
	assert.NoError(t, err)

	// 模拟发送登录消息
	loginMessage := constructLoginMessage() // 请根据你的协议构建登录消息
	sendMessage(t, conn, loginMessage)

	// 模拟接收登录响应
	receivedPacket := receiveMessage(t, conn)
	// 解码网络包
	decodedPacket, err := msg_packet.DecodeNetworkPacket(receivedPacket)
	if err != nil {
		fmt.Println("Error decoding network packet:", err)
		return
	}
	// 输出解码后的消息头和消息体
	fmt.Printf("Received Header: %+v\n", decodedPacket.Header)
	fmt.Printf("Received Body: %s\n", decodedPacket.Body)

	return
	// TODO: 解析登录响应并进行断言检查
	// ------------------------------------------------

	//// 模拟发送 ping 消息
	//pingMessage := constructPingMessage() // 请根据你的协议构建 ping 消息
	//sendMessage(t, conn, pingMessage)
	//// 模拟接收 ping 响应
	//pingResponse := receiveMessage(t, conn)
}

func constructLoginMessage() []byte {
	message := &cs.CMLogin{
		AccountId: "shyinyong",
		Password:  "123456",
	}
	messageBytes, _ := proto.Marshal(message)

	// 模拟创建一个网络包
	packet := msg_packet.NetworkPacket{
		Header: msg_packet.MessageHeader{
			PackageLen: 0, // 此处暂时设置为0，稍后会更新为正确的长度
			MsgID:      consts.CMMessageID_CMLogin,
			SeqID:      12345,
			MagicCode:  123,
			Reserved:   0,
		},
		Body: messageBytes,
	}
	packet.Header.PackageLen = uint16(consts.HeaderSize + len(packet.Body))
	packetBytes, err := msg_packet.EncodeNetworkPacket(packet)
	if err != nil {
		fmt.Println("Error encoding network packet:", err)
		return nil
	}
	return packetBytes

	// 根据你的协议构建登录消息的字节数组
	// 例如：message := &cs.CMLogin{...}
	// 然后使用 proto.Marshal 将消息序列化为字节数组
	// messageBytes, _ := proto.Marshal(message)
	// return messageBytes
}

func constructPingMessage() []byte {
	// 根据你的协议构建 ping 消息的字节数组
	// 例如：message := &cs.CMPing{...}
	// 然后使用 proto.Marshal 将消息序列化为字节数组
	// messageBytes, _ := proto.Marshal(message)
	// return messageBytes
	return nil
}

func sendMessage(t *testing.T, conn net.Conn, message []byte) {
	_, err := conn.Write(message)
	if err != nil {
		t.Fatalf("Error sending message: %v", err)
	}
}

func receiveMessage(t *testing.T, conn net.Conn) []byte {
	buffer := make([]byte, 1024) // 适当设置接收缓冲区的大小
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	return buffer[:n]
}
