package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	mutex sync.Mutex
}

func Start() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Print("Server started, waiting for connections...")

	var wg sync.WaitGroup
	server := &Server{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("Accepted connection to %v from %v\n", conn.LocalAddr(), conn.RemoteAddr())

		//go server.handleRequest(conn)
		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			defer conn.Close()
			server.handleConnection(conn)
		}(conn)
	}

	wg.Wait()
}

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

	// 处理请求消息
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// TODO: 处理业务逻辑
	// Unmarshal the protobuf message
	//handler.HandleMsg(conn, msgID)

	// 构造响应消息
	// 发送响应消息
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
