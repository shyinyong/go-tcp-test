package chat

import (
	"bufio"
	"fmt"
	"github.com/shyinyong/go-tcp-test/utils"
	"log"
	"net"
	"sync"
)

type Server struct {
	mutex         sync.Mutex
	rooms         map[string]*ChatRoom
	defaultRoom   *ChatRoom // Default chat room
	systemMessage chan string
	listener      net.Listener // Listener instance
}

func NewServer() *Server {
	return &Server{
		rooms:         make(map[string]*ChatRoom),
		defaultRoom:   GetDefaultChatRoom(),
		systemMessage: make(chan string),
	}
}

func (s *Server) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error starting server:", err)
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
	user := NewUser(conn)
	user.Username = username
	user.Conn = conn
	user.Writer = bufio.NewWriter(conn)

	s.defaultRoom.AddUser(user)

	room := GetDefaultChatRoom()
	onlineUserCount := len(room.Users)
	fmt.Printf("New join default room , current user count:%d\n", onlineUserCount)

	// Send a welcome message to the user
	//room.SendSystemMessage(fmt.Sprintf("Welcome, %s!", user.Username))

	// Handle user messages in a separate goroutine
	go user.listenForMessages()

	// Handle user disconnect
	defer user.disconnect()
}

func (s *Server) authenticateUser(conn net.Conn) (string, bool) {
	// Implement user authentication logic here
	// Read user credentials from the connection
	// Query the database or use any other method to authenticate the user
	// Return the username and a boolean indicating whether authentication succeeded
	username := utils.RandomUsername()
	return username, true // Replace with actual values
}

// 发送系统消息
func (s *Server) SendSystemMessage(message string) {
	s.systemMessage <- message
}

// 处理系统消息
func (s *Server) handleSystemMessages() {
	for {
		message := <-s.systemMessage
		s.defaultRoom.BroadcastSystemMessage(message)
	}
}
func (s *Server) broadcastSystemMessage(message string) {
	// Broadcast the system message to all users in the default room
	s.defaultRoom.BroadcastSystemMessage(message)
}

func (s *Server) Stop() {
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			log.Println("Error closing listener:", err)
		}
	}
}