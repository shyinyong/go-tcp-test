package utils

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/shyinyong/go-tcp-test/pb/chat"
	"testing"
)

func TestMessage(t *testing.T) {
	// 压包
	chatMessage := &chat.LoginReq{
		Username: "Alice",
		Password: "Bob",
	}
	serializedMessage, err := proto.Marshal(chatMessage)
	if err != nil {
		t.Errorf("error:%s", err)

	}
	msgType := 3
	msgID := 5
	fmt.Printf("body len：%d \n", len(serializedMessage))

	packed := PackMessage(uint16(msgType), uint32(msgID), serializedMessage)
	fmt.Printf("packed data:%v, size:%d\n", packed, len(packed))

	//var unpackedMessageType *Packet
	var unpackedMessage *Packet
	unpackedMessage, err = UnpackMessage(packed)
	if err != nil {
		t.Errorf("error:%s", err)
	}
	fmt.Println("MessageType:", unpackedMessage.MessageType)
	fmt.Println("MessageID:", unpackedMessage.MessageID)
	fmt.Println("MessageBody:", unpackedMessage.Body)

	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//chatMessage := &chat.LoginReq{}
	//err = proto.Unmarshal(unpackedMessage, chatMessage)
	//if err != nil {
	//	//fmt.Errorf("data error %v", err)
	//	return
	//}
	//
	//fmt.Println("Message Type:", unpackedMessageType)
	//fmt.Println("Message ID:", unpackedID)
	//fmt.Println("Sender:", chatMessage.Username)
	//fmt.Println("Content:", chatMessage.Password)

	//
	//// 压包
	//packedData, err := packMessage(MessageTypeA, 1, message)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("packed data:%v, size:%d\n", packedData, len(packedData))
	//
	//// 解包
	//unpackedMessageType, unpackedID, unpackedMessage, err := unpackMessage(packedData)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//chatMessage := &chat.LoginReq{}
	//err = proto.Unmarshal(unpackedMessage, chatMessage)
	//if err != nil {
	//	//fmt.Errorf("data error %v", err)
	//	return
	//}
	//
	//fmt.Println("Message Type:", unpackedMessageType)
	//fmt.Println("Message ID:", unpackedID)
	//fmt.Println("Sender:", chatMessage.Username)
	//fmt.Println("Content:", chatMessage.Password)
}
