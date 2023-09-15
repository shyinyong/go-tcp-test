package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	HeaderSize = 10 // 包头总长度
)

type Packet struct {
	msgType uint16
	msgId   uint32
	Body    []byte
}

// header size include:  msg type 2 + msg id 4 + msg body 4

func (p *Packet) GetBody() []byte {
	return p.Body
}

func (p *Packet) GetBodyLen() int {
	return len(p.Body)
}

func PackMessage(msgType uint16, msgId uint32, body []byte) []byte {
	packet := Packet{
		msgType: msgType,
		msgId:   msgId,
		Body:    body,
	}

	// 将 Packet 转换为二进制数据
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, packet.msgType)
	binary.Write(&buf, binary.BigEndian, packet.msgId)
	binary.Write(&buf, binary.BigEndian, uint32(len(packet.Body)))
	buf.Write(packet.Body)

	return buf.Bytes()
}

func UnpackMessage(data []byte) (*Packet, error) {
	if len(data) < HeaderSize {
		return nil, errors.New("insufficient data for packet header")
	}

	msgType := binary.BigEndian.Uint16(data[:2])
	msgId := binary.BigEndian.Uint32(data[2:6])
	bodyLength := binary.BigEndian.Uint32(data[6:10])

	if len(data) < int(HeaderSize+bodyLength) {
		return nil, errors.New("insufficient data for complete packet")
	}

	packet := &Packet{
		msgType: msgType,
		msgId:   msgId,
		Body:    data[HeaderSize : HeaderSize+bodyLength],
	}

	return packet, nil
}

//
//// msgTypes 消息类型定义
//type msgTypes uint8
//
//consts (
//	// msgTypeA 消息类型A
//	msgTypeA msgTypes = 1
//	// msgTypeB 消息类型B
//	msgTypeB msgTypes = 2
//	// msgTypeC 消息类型C
//	msgTypeC msgTypes = 3
//)
//
//// 压包方法
//func packMessage(msgType msgTypes, msgId uint32, chatMessage proto.Message) ([]byte, error) {
//	// 序列化聊天消息
//	serializedMessage, err := proto.Marshal(chatMessage)
//	if err != nil {
//		return nil, fmt.Errorf("failed to serialize chat message: %v", err)
//	}
//
//	// 消息类型（Message Type）：使用一个字节来表示消息的类型。在Go语言中，一个字节的范围是0到255，足够表示不同的消息类型。
//	// 消息ID（Message ID）：使用4个字节（32位整数）来表示消息的唯一标识符。这样可以提供足够的范围来表示大量的消息ID。
//	// 包长度（Packet Length）：使用4个字节（32位整数）来表示整个包的长度，即包头和消息体的总长度。这样可以确保解包时能够正确地读取整个包的数据。
//	// 因此，包头总共需要占用9个字节（1字节的消息类型 + 4字节的消息ID + 4字节的包长度）。
//
//	// 计算包长度
//	packetLength := 1 + 4 + 4 + len(serializedMessage)
//
//	// 构建包头
//	packetHeader := make([]byte, 9)
//	packetHeader[0] = byte(msgType)                                // 消息类型
//	binary.BigEndian.PutUint32(packetHeader[1:], msgId)            // 消息ID
//	binary.BigEndian.PutUint32(packetHeader[5:], uint32(packetLength)) // 包长度
//
//	// 拼接包头和消息体
//	packetData := append(packetHeader, serializedMessage...)
//
//	// 返回压包后的数据
//	return packetData, nil
//}
//
//// 解包方法
//func unpackMessage(packedData []byte) (msgTypes, uint32, []byte, error) {
//	fmt.Printf("unpackMessage, packeddata size:%d \n", len(packedData))
//
//	if len(packedData) < 9 {
//		return 0, 0, nil, fmt.Errorf("invalid packed data: insufficient length")
//	}
//
//	// 解析包头
//	msgType := msgTypes(packedData[0])               // 消息类型
//	msgId := binary.BigEndian.Uint32(packedData[1:5])    // 消息ID
//	packetLength := binary.BigEndian.Uint32(packedData[5:9]) // 包长度
//
//	fmt.Printf("msgType :%d \n", msgType)
//	fmt.Printf("msgId :%d \n", msgId)
//	fmt.Printf("packetLength :%d \n", packetLength)
//
//	if uint32(len(packedData)) < packetLength {
//		return 0, 0, nil, fmt.Errorf("invalid packed data: incorrect length")
//	}
//
//	// 返回解包后的消息类型、消息ID和消息体
//	return msgType, msgId, packedData[9 : 9+packetLength], nil
//}
