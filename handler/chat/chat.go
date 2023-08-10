package chat

import (
	"bufio"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/shyinyong/go-tcp-test/pb/chat"
	"github.com/shyinyong/go-tcp-test/utils"
	"log"
	"net"
	"sync"
)

type Server struct {
	mutex    sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	// ...
	rooms         map[string]*ChatRoom
	defaultRoom   *ChatRoom // Default chat room `Only server send message to me`
	worldRoom     *ChatRoom // World chat room
	userRooms     map[*User]*ChatRoom
	systemMessage chan string
}

func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		ctx:    ctx,
		cancel: cancel,
		// ...
		rooms:         make(map[string]*ChatRoom),
		defaultRoom:   GetDefaultChatRoom(),
		worldRoom:     NewChatRoom("World Chat"),
		systemMessage: make(chan string),
	}
}

func (s *Server) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// TODO 尽量避免在代码中直接使用 log.Fatal，因为这会导致整个程序的终止。你可以考虑将错误返回到调用函数，然后让调用者决定如何处理这个错误。
		log.Println("Error starting server:", err)
		return
	}
	s.listener = listener // Save the listener instance
	fmt.Println("Chat server started. Listening on", address)

	// 启动系统消息处理
	go s.handleSystemMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go s.handleChatConnection(conn)
	}
}

func (s *Server) handleChatConnection(conn net.Conn) {
	defer conn.Close()

	// Perform user authentication
	username, authenticated := s.authenticateUser(conn)
	if !authenticated {
		log.Println("Authentication failed for connection:", conn.RemoteAddr())
		return
	}

	// Create a User instance
	user := NewUser(conn, s)
	user.Username = username
	user.Conn = conn
	user.Writer = bufio.NewWriter(conn)

	// Handle user disconnect
	defer user.disconnect()

	s.defaultRoom.AddUser(user)
	s.worldRoom.AddUser(user)

	// Send a welcome message to the user
	s.defaultRoom.BroadcastSystemMessage(fmt.Sprintf("Welcome, %s!", user.Username))
	// Handle user messages in a separate goroutine
	go user.listenForMessages()
}

func (s *Server) authenticateUser(conn net.Conn) (string, bool) {
	// Implement user authentication logic here
	// Read user credentials from the connection
	// Query the database or use any other method to authenticate the user
	// Return the username and a boolean indicating whether authentication succeeded
	username := utils.RandomUsername()
	return username, true // Replace with actual values
}

// SendSystemMessage 发送系统消息
func (s *Server) SendSystemMessage(message string) {
	s.systemMessage <- message
}

// 处理系统消息
func (s *Server) handleSystemMessages() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case message := <-s.systemMessage:
			s.worldRoom.BroadcastSystemMessage(message)
		}
	}
}

func (s *Server) Stop() {
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			log.Println("Error closing listener:", err)
		}
	}
	s.cancel() // cancel all goroutines
}

// 发送世界聊天消息
func (s *Server) SendWorldChatMessage(sender *User, content string) {
	chatMessage := &chat.ChatMessage{
		Username: sender.Username,
		Content:  content,
	}

	// 序列化消息并广播给世界聊天室的所有用户
	serializedMessage, err := proto.Marshal(chatMessage)
	if err != nil {
		log.Println("Error serializing message:", err)
		return
	}

	s.worldRoom.BroadcastMessage(serializedMessage)
}

// 处理接收到的世界聊天消息
func (s *Server) handleWorldChatMessage(sender *User, serializedMessage []byte) {
	chatMessage := &chat.ChatMessage{}
	err := proto.Unmarshal(serializedMessage, chatMessage)
	if err != nil {
		log.Println("Error deserializing message:", err)
		return
	}

	// 在世界聊天室广播消息
	s.worldRoom.Broadcast(sender, fmt.Sprintf("[%s] %s", chatMessage.Username, chatMessage.Content))
}
