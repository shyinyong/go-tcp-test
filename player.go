package main

import (
	"fmt"
	"github.com/shyinyong/go-tcp-test/pb/cs"
	"google.golang.org/protobuf/proto"
)

func CMPing(message []byte) []byte {
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
	return nil
}

// CMLogin 处理函数
func CMLogin(message []byte) []byte {
	account := "ABC"
	fmt.Println("Received SMLogin message")

	rspMsg := &cs.SMLogin{}
	rspMsg.AccountId = &account
	rspBytes, _ := proto.Marshal(rspMsg)
	return rspBytes

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
	return nil
}
