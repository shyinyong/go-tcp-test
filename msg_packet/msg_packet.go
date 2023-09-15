package msg_packet

import (
	"encoding/binary"
	"fmt"
	"github.com/shyinyong/go-tcp-test/consts"
	"net"
)

type WspCliConn struct {
	Conn         net.Conn
	SocketHandle uint16
}

type MsgHdr struct {
	MsgId        uint16
	SeqId        uint32
	SocketHandle uint16
	IpSaddr      uint32

	Conn    net.Conn
	Data    []byte
	Msg     interface{}
	Context interface{}
}

type NetMsgHandler[T any] struct {
	MsgId     int
	HandlerId int
	ParseCb   func([]byte) interface{}
	Cb        func(*MsgHdr, T)
}

type NetMsg interface {
	GetNetMsgId() uint16
}

type WSProxyPackHead_C struct {
	PackLen      uint16
	MsgId        uint16
	SeqId        uint32
	MagicCode    uint16
	SocketHandle uint16

	IpSaddr uint64
}

type WSProxyPackHead_S struct {
	PackLen    uint16
	MsgId      uint16
	SeqId      uint32
	MagicCode  uint16
	RpcErrCode uint16

	SocketHandle uint16
	ExtLen       uint16
}

func (this *WspCliConn) IsValid() bool {
	return this.Conn != nil && this.SocketHandle != 0
}

func (this *WspCliConn) Reset() {
	this.Conn = nil
	this.SocketHandle = 0
}

func (this *MsgHdr) GetSocket() WspCliConn {
	socket := WspCliConn{}
	socket.Conn = this.Conn
	socket.SocketHandle = this.SocketHandle
	return socket
}

// -----------------------------------

type MessageHeader struct {
	PackageLen uint16
	MsgID      uint16
	SeqID      uint32
	MagicCode  uint16
	Reserved   uint16
}

type NetworkPacket struct {
	Header MessageHeader
	Body   []byte
}

// EncodeNetworkPacket 编码网络包
func EncodeNetworkPacket(packet NetworkPacket) ([]byte, error) {
	headerBytes := make([]byte, consts.HeaderSize)
	binary.BigEndian.PutUint16(headerBytes[:2], packet.Header.PackageLen)
	binary.BigEndian.PutUint16(headerBytes[2:4], packet.Header.MsgID)
	binary.BigEndian.PutUint32(headerBytes[4:8], packet.Header.SeqID)
	binary.BigEndian.PutUint16(headerBytes[8:10], packet.Header.MagicCode)
	binary.BigEndian.PutUint16(headerBytes[10:12], packet.Header.Reserved)

	return append(headerBytes, packet.Body...), nil
}

// DecodeNetworkPacket 解码网络包
func DecodeNetworkPacket(data []byte) (NetworkPacket, error) {
	if len(data) < consts.HeaderSize {
		return NetworkPacket{}, fmt.Errorf("insufficient data for header")
	}

	header := MessageHeader{
		PackageLen: binary.BigEndian.Uint16(data[:2]),
		MsgID:      binary.BigEndian.Uint16(data[2:4]),
		SeqID:      binary.BigEndian.Uint32(data[4:8]),
		MagicCode:  binary.BigEndian.Uint16(data[8:10]),
		Reserved:   binary.BigEndian.Uint16(data[10:12]),
	}

	body := data[consts.HeaderSize:]

	return NetworkPacket{
		Header: header,
		Body:   body,
	}, nil
}
