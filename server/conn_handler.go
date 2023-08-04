package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/shyinyong/go-tcp-test/handler"
	"github.com/shyinyong/go-tcp-test/pb/address"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
)

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()
	defer s.handleDisconnection(conn)

	reqData, err := readRequest(conn)
	if err != nil {
		// 处理错误
		return
	}

	// 解析请求消息
	reqMsg := &address.Person{}
	err = proto.Unmarshal(reqData, reqMsg)
	if err != nil {
		// 处理错误
		return
	}

	// 处理请求消息
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// TODO: 处理业务逻辑

	// 构造响应消息
	respMsg := &address.Person{}
	respData, err := proto.Marshal(respMsg)
	if err != nil {
		// 处理错误
		return
	}

	// 发送响应消息
	err = writeResponse(conn, respData)
	if err != nil {
		// 处理错误
		return
	}
}

func (s *Server) handleDisconnection(conn net.Conn) {

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	var requestData bytes.Buffer
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			// Handle error
			break
		}
		requestData.Write(buf[:n])
		if n < len(buf) {
			break
		}
	}
	handler.HandleMsg(conn, requestData.Bytes())
}

func readRequest(conn net.Conn) ([]byte, error) {
	// 读取消息头中的长度字段
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lenBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf)

	// 读取消息体
	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil

	//var requestData = bytes.Buffer{}
	//buf := make([]byte, 256)
	//for {
	//	n, err := conn.Read(buf)
	//	if err != nil {
	//		if err != io.EOF {
	//			fmt.Println("read error:", err)
	//		}
	//		// Handle error
	//		break
	//	}
	//	requestData.Write(buf[:n])
	//	if n < len(buf) {
	//		break
	//	}
	//}
	//return requestData.Bytes(), nil
}

func writeResponse(conn net.Conn, data []byte) error {
	// 构造消息头
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))

	// 发送消息头和消息体
	_, err := conn.Write(append(header, data...))
	return err

	// step1 https://protobuf.dev/getting-started/gotutorial/
	// step2 go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	// step3 protoc -I=D:\projects\go-tcp-test --go_out=D:\projects\go-tcp-test\pb D:\projects\go-tcp-test\proto\*.proto

	// 在读取请求消息的逻辑中，首先读取消息头中的长度字段，然后根据长度字段读取消息体。
	// 在发送响应消息的逻辑中，先构造消息头，然后将消息头和消息体合并后一次性发送。
	// 注意在发送消息时，需要将消息头和消息体合并为一个字节数组进行发送。
}
