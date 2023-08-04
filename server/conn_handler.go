package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/shyinyong/go-tcp-test/handler"
	"github.com/shyinyong/go-tcp-test/pb/request"
	"github.com/shyinyong/go-tcp-test/pb/response"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"time"
)

const (
	HeartbeatInterval = 10 * time.Second // 心跳间隔时间
	TimeoutDuration   = 30 * time.Second // 连接超时时间
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer s.handleDisconnection(conn)

	// 设置超时时间
	conn.SetDeadline(time.Now().Add(TimeoutDuration))
	// 创建心跳计时器
	heartbeatTicker := time.NewTicker(HeartbeatInterval)
	defer heartbeatTicker.Stop()

	// 读取和处理心跳消息
	go func() {
		for {
			select {
			case <-heartbeatTicker.C:
				_, err := conn.Write([]byte("heartbeat"))
				if err != nil {
					fmt.Println("Error sending heartbeat:", err)
					return
				}
			}
		}
	}()

	reqData, err := readRequest(conn)
	if err != nil {
		// 处理错误
		return
	}

	// 解析请求消息
	reqMsg := &request.Request{}
	err = proto.Unmarshal(reqData, reqMsg)
	if err != nil {
		// 处理错误
		return
	}

	// 处理请求消息
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// TODO: 处理业务逻辑
	handler.HandleMsg(conn, reqMsg)

	// 构造响应消息
	respMsg := &response.Response{}
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

// 添加心跳响应处理逻辑
//func handleHeartbeatResponse(c net.Conn, uid int32, onlinePlayers *OnlinePlayers) {
//	onlinePlayers.Lock()
//	defer onlinePlayers.Unlock()
//
//	// 更新最近心跳时间
//	onlinePlayers.lastHeartbeats[uid] = time.Now()
//}

// 处理断线时移除玩家信息
func (s *Server) handleDisconnection(conn net.Conn) {
	//onlinePlayers.Lock()
	//delete(onlinePlayers.players, uid)
	//delete(onlinePlayers.lastHeartbeats, uid)
	//onlinePlayers.Unlock()
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

	request := &request.Request{}
	request.Msg = requestData.String()
	handler.HandleMsg(conn, request)
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
}
