package chat

import (
	"bufio"
	"context"
	"fmt"
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
	systemMessage chan *chat.SystemMessage
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
		systemMessage: make(chan *chat.SystemMessage),
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

func (s *Server) Stop() {
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			log.Println("Error closing listener:", err)
		}
	}
	s.cancel() // cancel all goroutines
}

func (s *Server) handleChatConnection(conn net.Conn) {
	defer conn.Close()

	// Perform user authentication
	username, authenticated := s.authenticateUser(conn)
	if !authenticated {
		log.Println("Authentication failed for connection:", conn.RemoteAddr())
		return
	}

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
	// Query the database or use any other method to authenticate the user
	// Return the username and a boolean indicating whether authentication succeeded
	username := utils.RandomUsername()
	return username, true // Replace with actual values
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

// 发送世界聊天消息
func (s *Server) SendWorldChatMessage(sender *User, content string) {
	// 构建世界聊天消息
	chatMessage := &chat.ChatMessage{
		SenderUsername: sender.Username,
		Content:        content,
		ChatType:       chat.ChatMessage_ChatType_WORLD,
	}

	// 发送世界聊天消息
	s.worldRoom.BroadcastMessage(chatMessage)
}
func (s *Server) SendGuildChatMessage(sender *User, content string) {
	// 构建公会聊天消息
	chatMessage := &chat.ChatMessage{
		SenderUsername: sender.Username,
		Content:        content,
		ChatType:       chat.ChatMessage_ChatType_GUILD,
	}

	// 获取用户所在的公会聊天频道（根据您的逻辑实现）
	guildChatRoom := sender.Guild.ChatRoom
	// 发送公会聊天消息
	guildChatRoom.BroadcastMessage(chatMessage)

	// 发送公会聊天消息
	// 这里需要根据用户所在的公会找到对应的公会聊天频道进行消息分发
	// 暂时留空
}

func (s *Server) SendPrivateChatMessage(sender *User, receiver *User, content string) {
	// 构建私人聊天消息
	chatMessage := &chat.ChatMessage{
		SenderUsername:   sender.Username,
		Content:          content,
		ChatType:         chat.ChatMessage_ChatType_PRIVATE,
		ReceiverUsername: receiver.Username,
	}

	// 发送私人聊天消息
	// 这里需要判断目标用户是否在线，如果在线则直接发送，如果不在线则暂存或发送离线通知
	// 查找目标用户
	targetUser := s.FindUserByUsername(receiver.Username)
	if targetUser != nil {
		// 发送私人聊天消息给目标用户
		targetUser.PrivateMessageChannel <- chatMessage
	} else {
		log.Printf("User %s not found or not in the same room.", receiver.Username)
	}

}
