package main

import (
	"encoding/binary"
	"fmt"
	protobuf_cs "github.com/shyinyong/go-tcp-test/pb/cs"
	"google.golang.org/protobuf/proto"
	"net"
)

const (
	msgHeadSize = 12
	magicCode   = 0x4b53 // KS的ASCII码
)

type MsgHead struct {
	PackageLen uint16 // 双字节网络包长度
	MsgID      uint16 // 双字节消息id
	SeqID      uint32 // 4字节序号id
	MagicCode  uint16 // 2字节魔数，并且为固定常数KS，占位符客户端不需什么处理
	Reserved   uint16 // 保留
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Error connecting", err)
		return
	}
	defer conn.Close() // Ensure the connection is closed after the test finishes

	// 构造消息头
	msgHead := MsgHead{
		PackageLen: msgHeadSize,
		MsgID:      1,
		SeqID:      1,
		MagicCode:  magicCode,
		Reserved:   0,
	}
	// 序列化消息头
	msgHeadData := make([]byte, msgHeadSize)
	binary.BigEndian.PutUint16(msgHeadData[:2], msgHead.PackageLen)
	binary.BigEndian.PutUint16(msgHeadData[2:4], msgHead.MsgID)
	binary.BigEndian.PutUint32(msgHeadData[4:8], msgHead.SeqID)
	binary.BigEndian.PutUint16(msgHeadData[8:10], msgHead.MagicCode)
	binary.BigEndian.PutUint16(msgHeadData[10:12], msgHead.Reserved)

	// 构造消息体（protobuf数据）
	accountId := "123456"
	var ProtoVersion int32 = 5
	msgBody := &protobuf_cs.CMListFriend{
		AccountId:    &accountId,
		ProtoVersion: &ProtoVersion,
		SessionId:    &accountId,
	}
	protoData, err := proto.Marshal(msgBody)
	if err != nil {
		fmt.Println("序列化protobuf数据失败:", err)
		return
	}

	// 构造完整的网络包
	message := append(msgHeadData, protoData...)

	// 发送消息
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("发送消息失败:", err)
		return
	}

	fmt.Println("消息发送成功")
}
