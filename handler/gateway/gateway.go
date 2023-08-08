package gateway

import (
	"fmt"
	"github.com/shyinyong/go-tcp-test/handler"
	protobuf "github.com/shyinyong/go-tcp-test/pb/message"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	mutex sync.Mutex

	heartbeatInterval time.Duration
	heartbeatTimer    *time.Timer
	serverAddresses   map[string]string
}

const heartbeatInterval = 10

func NewServer() *Server {
	serverAddresses := map[string]string{
		"login": "localhost:8081",
		//"game":  "game-server:12347",
		//"chat":  "chat-server:12348",
		// Add more server addresses if needed...
	}

	return &Server{
		// 初始化代码...
		heartbeatInterval: 10 * time.Second,
		heartbeatTimer:    time.NewTimer(heartbeatInterval),
		serverAddresses:   serverAddresses,
		//s.loginServerAddr = "localhost:8081"
		//s.gameServerAddr = "localhost:12347"
		//s.chatServerAddr = "localhost:12348"
	}
}

// SetServerAddresses sets the addresses of different servers.
func (s *Server) SetServerAddresses(addresses map[string]string) {
	s.serverAddresses = addresses
}

func (s *Server) AddServerAddr(messageType protobuf.ServerType, address string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//s.serverAddrs[messageType] = address
}

func (s *Server) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()
	fmt.Println("Gateway server started. Listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go s.handleGatewayConnection(conn)
	}
}

func (s *Server) handleGatewayConnection(conn net.Conn) {
	defer conn.Close()
	for {
		fmt.Println("new connection")

		// 读取客户端消息
		data, err := s.readMessage(conn)
		if err != nil {
			log.Println("Error reading client message:", err)
			return
		}

		// 反序列化客户端消息
		clientMsg := &protobuf.ClientMessage{}
		err = proto.Unmarshal(data, clientMsg)
		if err != nil {
			log.Println("Error unmarshaling client message:", err)
			return
		}

		// 处理心跳请求
		if clientMsg.Type == protobuf.ClientMessage_HEARTBEAT_REQUEST {
			handler.HandleHeartbeat(conn, clientMsg, time.Now())
			continue
		}

		fmt.Printf("remote server adds:%v \n\n", s.serverAddresses)

		// 转发消息到相应的服务器
		switch clientMsg.Type {
		case protobuf.ClientMessage_LOGIN_REQUEST,
			protobuf.ClientMessage_RELOGIN_REQUEST,
			protobuf.ClientMessage_LOGOUT_REQUEST:
			// 转发登录、重新登录和退出请求到登录服务器

			remoteServerAdd := s.serverAddresses["login"]
			fmt.Printf("remote server add:%s", remoteServerAdd)
			fmt.Printf("end")

			s.forwardToServer(remoteServerAdd, clientMsg)
		case protobuf.ClientMessage_ENTER_GAME_REQUEST,
			protobuf.ClientMessage_BATTLE_START_REQUEST,
			protobuf.ClientMessage_BATTLE_FAIL_REQUEST:
			// 转发进入游戏、战斗开始和战斗失败请求到游戏服务器
			s.forwardToServer(s.serverAddresses["game"], clientMsg)
		case protobuf.ClientMessage_SEND_TEXT_MESSAGE_REQUEST:
			// 转发发送文字消息请求到聊天服务器
			s.forwardToServer(s.serverAddresses["chat"], clientMsg)

		default:
			log.Println("Unknown message type")
		}
	}
}

func (s *Server) forwardToServer(serverAddr string, clientMsg *protobuf.ClientMessage) {
	// Check if the server address exists in the serverAddresses map
	if _, ok := s.serverAddresses[serverAddr]; !ok {
		log.Printf("Server address '%s' not found in serverAddresses map", serverAddr)
		return
	}

	// Send the client message to the specified server
	err := s.sendMessageToServer(serverAddr, clientMsg)
	if err != nil {
		log.Printf("Error forwarding message to server '%s': %v", serverAddr, err)
		return
	}
}

// sendMessageToServer sends the client message to the specified server.
func (s *Server) sendMessageToServer(serverAddr string, clientMsg *protobuf.ClientMessage) error {
	// Get the connection to the specified server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Serialize the client message
	data, err := proto.Marshal(clientMsg)
	if err != nil {
		return err
	}

	// Send the client message to the server
	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) readMessage(conn net.Conn) ([]byte, error) {
	// 创建一个缓冲区用于读取客户端发送的消息
	buffer := make([]byte, 1024)

	// 从连接中读取数据
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	// 提取有效数据
	data := buffer[:n]

	return data, nil
}

//func (s *Server) ConnectToServers() error {
//	s.mutex.Lock()
//	defer s.mutex.Unlock()
//
//	for messageType, address := range s.serverAddrs {
//		conn, err := net.Dial("tcp", address)
//		if err != nil {
//			log.Printf("Error connecting to %s server: %v\n", messageType, err)
//			return err
//		}
//		s.serverConns[messageType] = conn
//	}
//
//	return nil
//}
