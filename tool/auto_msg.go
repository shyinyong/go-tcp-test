package main

import (
	"fmt"
	"os"
	"strings"
)

const autoMsgTemplate = `package cs

import (
	"google.golang.org/protobuf/proto"
	"github.com/shyinyong/go-tcp-test/consts"
	"github.com/shyinyong/go-tcp-test/msg_packet"
)

type CsNetMsgHandler msg_packet.NetMsgHandler[MsgHandler];

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
%s
}

%s

%s

`

func generateAutoMsgFile(messageIDs []string) {
	// 生成函数签名列表
	functionSignatures := []string{}
	implFunctions := []string{}
	generateGetNetMsgIds := []string{}
	for _, messageID := range messageIDs {
		functionSignatures = append(functionSignatures, generateFunctionSignature(messageID))
		implFunctions = append(implFunctions, generateImplFunction(messageID))
		generateGetNetMsgIds = append(generateGetNetMsgIds, generateGetNetMsgId(messageID))
	}

	// 构建 auto_msg.go 文件内容
	autoMsgContent := fmt.Sprintf(
		autoMsgTemplate,
		strings.Join(functionSignatures, "\n"),
		strings.Join(implFunctions, "\n"),
		strings.Join(generateGetNetMsgIds, "\n"),
	)

	// 将内容写入 auto_msg.go 文件
	wordDir, _ := os.Getwd()
	err := os.WriteFile(wordDir+"/pb/cs/auto_msg.go", []byte(autoMsgContent), os.ModePerm)
	if err != nil {
		fmt.Println("Error writing auto_msg.go:", err)
	}
}

func generateFunctionSignature(messageID string) string {
	return fmt.Sprintf("%s(hdr *msg_packet.MsgHdr, msg *%s)", messageID, messageID)
}

func generateImplFunction(messageID string) string {
	return fmt.Sprintf(`
func (msgI *MsgHandlerImpl) %s(hdr *msg_packet.MsgHdr, msg *%s) {
	// 实现 %s 处理逻辑
}
`, messageID, messageID, messageID)
}

func generateGetNetMsgId(messageID string) string {
	return fmt.Sprintf(`
func (cm *%s) GetNetMsgId() uint16 {
	return uint16(consts.CMMessageID_%s)
}

`, messageID, messageID)
}

func main() {
	// 定义要生成的消息ID列表
	messageIDs := []string{
		"CMPing",
		"CMLogin",
		// 添加其他消息ID
	}

	// 生成 auto_msg.go 文件
	generateAutoMsgFile(messageIDs)
	fmt.Println("auto_msg.go generated successfully.")
}
