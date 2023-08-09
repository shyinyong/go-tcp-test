package gateway

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shyinyong/go-tcp-test/config"
	"github.com/shyinyong/go-tcp-test/handler"
	protobuf "github.com/shyinyong/go-tcp-test/pb/message"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	mutex  sync.Mutex
	config config.Config
	store  *sqlx.DB

	heartbeatInterval time.Duration
	heartbeatTimer    *time.Timer
	serverAddresses   map[string]string
}

const heartbeatInterval = 10

func NewServer(cfg config.Config, store *sqlx.DB) *Server {
	return &Server{
		// 初始化代码...
		config:            cfg,
		store:             store,
		heartbeatInterval: 10 * time.Second,
		heartbeatTimer:    time.NewTimer(heartbeatInterval),
	}
}

// SetServerAddresses sets the addresses of different servers.
func (s *Server) SetServerAddresses(addresses map[string]string) {
	s.serverAddresses = addresses
}

func (s *Server) AddServerAddr(messageType protobuf.ServerType, address string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
}

func (s *Server) Start(addresses []string) error {
	var wg sync.WaitGroup
	wg.Add(len(addresses))

	for _, address := range addresses {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return fmt.Errorf("error starting listener on address %s: %w", address, err)
		}

		log.Printf("Gateway server started. Listening on:%s \n", address)

		go func(l net.Listener) {
			defer l.Close()
			defer wg.Done()
			s.acceptConnections(l)
		}(listener)
	}

	wg.Wait()
	return nil
}

func (s *Server) acceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting client connection:", err)
			continue
		}
		fmt.Printf("conn:%s", conn.RemoteAddr())

		// Create a context for this connection
		ctx, cancel := context.WithCancel(context.Background())

		go s.handleGatewayConnection(ctx, conn, true)

		// Listen for the context to be canceled and close the connection
		go func() {
			<-ctx.Done()
			conn.Close()
		}()

		// When the connection is done, cancel the context
		defer cancel()
	}
}

func (s *Server) handleGatewayConnection(ctx context.Context, conn net.Conn, shouldClose bool) {
	// Only close the connection if the shouldClose flag is true
	if shouldClose {
		defer conn.Close()
	}
	for {
		select {
		case <-ctx.Done():
			// Context canceled, stop handling the connection
			return
		default:
			data, err := s.readMessage(conn)
			if err != nil {
				if err != io.EOF {
					log.Println("Error reading client message:", err)
				}
				//log.Println("Error reading client message:", err)
				return
			}

			clientMsg := &protobuf.ClientMessage{}
			err = proto.Unmarshal(data, clientMsg)
			if err != nil {
				log.Println("Error unmarshaling client message:", err)
				return
			}

			// 处理心跳请求
			if clientMsg.Type == protobuf.ClientMessage_HEARTBEAT_REQUEST {
				//s.handleHeartbeatRequest(conn, clientMsg)
				handler.HandleHeartbeat(conn, clientMsg, time.Now())
				continue
			}

			// 转发消息到相应的服务器
			s.forwardMessage(conn, clientMsg)
		}
	}
}

func (s *Server) forwardMessage(conn net.Conn, clientMsg *protobuf.ClientMessage) {
	var serverAddr string

	switch clientMsg.Type {
	case protobuf.ClientMessage_LOGIN_REQUEST,
		protobuf.ClientMessage_RELOGIN_REQUEST,
		protobuf.ClientMessage_LOGOUT_REQUEST:
		// 转发登录、重新登录和退出请求到登录服务器
		serverAddr = s.serverAddresses["login"]
	case protobuf.ClientMessage_ENTER_GAME_REQUEST,
		protobuf.ClientMessage_BATTLE_START_REQUEST,
		protobuf.ClientMessage_BATTLE_FAIL_REQUEST:
		// 转发进入游戏、战斗开始和战斗失败请求到游戏服务器
		serverAddr = s.serverAddresses["game"]
	case protobuf.ClientMessage_SEND_TEXT_MESSAGE_REQUEST:
		// 转发发送文字消息请求到聊天服务器
		serverAddr = s.serverAddresses["chat"]
	default:
		log.Println("Unknown message type")
		return
	}

	// Get the connection to the specified server
	serverConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("Error connecting to %s server: %v", serverAddr, err)
		return
	}
	defer serverConn.Close()

	// Serialize the client message
	data, err := proto.Marshal(clientMsg)
	if err != nil {
		log.Printf("Error marshaling client message: %v", err)
		return
	}

	// Send the client message to the server
	_, err = serverConn.Write(data)
	if err != nil {
		log.Printf("Error sending message to %s server: %v", serverAddr, err)
		return
	}

	log.Printf("Forwarded message to %s server: %v \n", serverAddr, clientMsg)
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
