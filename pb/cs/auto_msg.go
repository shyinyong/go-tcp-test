package cs

import (
	"github.com/shyinyong/go-tcp-test/consts"
	"github.com/shyinyong/go-tcp-test/msg_packet"
)

type CsNetMsgHandler msg_packet.NetMsgHandler[MsgHandler]

type MsgHandlerImpl struct {
}

var handlers [2000]*CsNetMsgHandler

func GetNetMsgHandler(msgId uint16) *CsNetMsgHandler {
	handler := handlers[msgId]
	return handler
}

func DispatchMsg(handler *CsNetMsgHandler, hdr *msg_packet.MsgHdr, msgHandler MsgHandler) {
	handler.Cb(hdr, msgHandler)
}

func RegHandlerId(msgId int, handlerId int) {
	handler := handlers[msgId]
	handler.HandlerId = handlerId
}

func ParsePb(msgId uint16, data []byte) interface{} {
	handler := handlers[msgId]
	if handler == nil {
		return nil
	}
	return handler.ParseCb(data)
}

type MsgHandler interface {
	CMPing(hdr *msg_packet.MsgHdr, msg *CMPing)
	CMLogin(hdr *msg_packet.MsgHdr, msg *CMLogin)
}

func (msgI *MsgHandlerImpl) CMPing(hdr *msg_packet.MsgHdr, msg *CMPing) {
	// 实现 CMPing 处理逻辑
}

func (msgI *MsgHandlerImpl) CMLogin(hdr *msg_packet.MsgHdr, msg *CMLogin) {
	// 实现 CMLogin 处理逻辑
}

func (cm *CMPing) GetNetMsgId() uint16 {
	return uint16(consts.CMMessageID_CMPing)
}

func (cm *CMLogin) GetNetMsgId() uint16 {
	return uint16(consts.CMMessageID_CMLogin)
}
