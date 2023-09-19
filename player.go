package main

import (
	"fmt"
	"github.com/shyinyong/go-tcp-test/msg_packet"
	"github.com/shyinyong/go-tcp-test/pb/cs"
	"google.golang.org/protobuf/proto"
)

type Player struct {
	cs.MsgHandlerImpl
	socket    msg_packet.WspCliConn
	accountId string
}

// need the method: CMLogin(hdr *msg_packet.MsgHdr, msg *CMLogin)
// have the method: CMLogin(message []byte) []byte

func (p *Player) CMPing(hdr *msg_packet.MsgHdr, msg *cs.CMPing) {
	fmt.Println("Received SMPing message")
	// 在这里处理 CMPing 消息，返回 SMPing 消息或其他响应
	// 示例代码：
	// cmPing := &cs.CMPing{}
	// 解码消息，处理消息内容
	// ...
	// 返回 SMPing 消息
	// rspMsg := &cs.SMPing{}
	// 填充 rspMsg 内容
	// ...
	// 返回响应消息的字节数组
	// rspBytes, _ := proto.Marshal(rspMsg)
	// return rspBytes
}

// CMLogin 处理函数
func (p *Player) CMLogin(hdr *msg_packet.MsgHdr, msg *cs.CMLogin) {
	account := "ABC"
	fmt.Println("Received SMLogin message")

	rspMsg := &cs.SMLogin{}
	rspMsg.AccountId = &account
	p.SendMsg(rspMsg)

	// 在这里处理 CMLogin 消息，返回 SMLogin 消息或其他响应
	// 示例代码：
	// cmLogin := &cs.CMLogin{}
	// 解码消息，处理消息内容
	// ...
	// 返回 SMLogin 消息
	// rspMsg := &cs.SMLogin{}
	// 填充 rspMsg 内容
	// ...
	// 返回响应消息的字节数组
	// rspBytes, _ := proto.Marshal(rspMsg)
	// return rspBytes
}

func (p *Player) SendMsg(rspMsg proto.Message) {
	// msg_packet.wspListener.sendProxyMsg(p.socket.Conn, p.socket.SocketHandle, rspMsg)
}
