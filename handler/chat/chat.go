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
	rooms         map[string]*Room
	defaultRoom   *Room // Default chat room `Only server send message to me`
	worldRoom     *Room // World chat room
	userRooms     map[*User]*Room
	Users         map[int64]*User
	systemMessage chan *chat.SystemMessage
}

func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		ctx:    ctx,
		cancel: cancel,
		// ...
		rooms:         make(map[string]*Room),
		defaultRoom:   GetDefaultRoom(),
		worldRoom:     NewChatRoom("WorldRoom"),
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
	defer listener.Close()

	s.listener = listener // Save the listener instance
	fmt.Println("Chat server started. Listening on", address)

	// 启动系统消息处理
	go s.handleSystemMessages()

	// 使用协程池来管理连接处理的goroutine数量
	pool := NewPool(100) // 创建大小为100的协程池
	defer pool.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		// 提交连接处理任务到协程池
		pool.Submit(func() {
			go s.handleChatConnection(conn)
		})
	}
}

func (s *Server) Stop() {
	if s.listener != nil {
		_ = s.listener.Close() // Ignore error
	}
	s.cancel() // cancel all goroutines
}

func (s *Server) handleChatConnection(conn net.Conn) {
	log.Println("New connection")
	defer conn.Close()

	// 创建一个上下文，用于控制连接的生命周期
	connCtx, connCancel := context.WithCancel(s.ctx)
	defer connCancel() // 保证连接关闭时取消所有 goroutine

	// Perform user authentication
	username, authenticated := s.authenticateUser(conn)
	if !authenticated {
		log.Println("Authentication failed for connection:", conn.RemoteAddr())
		return
	}
	var userId = username + string(utils.RandomUserID())

	user := NewUser(conn, s)
	user.UserID = userId
	user.Username = username
	user.Conn = conn
	user.Writer = bufio.NewWriter(conn)
	go user.listenForMessages() // 启动消息监听协程

	// 处理用户断开连接
	defer user.disconnect()

	s.defaultRoom.AddUser(user)
	s.worldRoom.AddUser(user)

	// Send a welcome message to the user
	s.worldRoom.BroadcastSystemMessage(&chat.SystemMessage{
		Content: fmt.Sprintf("Welcome, %s!", user.Username),
	})

	// Handle user messages in a separate goroutine
	go user.listenForMessagesWithContext(connCtx)
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

	// 序列化系统消息
	messageBytes, err := proto.Marshal(chatMessage)
	if err != nil {
		log.Println("Error serializing system message:", err)
		return
	}
	// 发送世界聊天消息
	s.worldRoom.BroadcastMessage(messageBytes)
}
func (s *Server) SendGuildChatMessage(sender *User, content string) {
	// 构建公会聊天消息
	//chatMessage := &chat.ChatMessage{
	//	SenderUsername: sender.Username,
	//	Content:        content,
	//	ChatType:       chat.ChatMessage_ChatType_GUILD,
	//}

	// 获取用户所在的公会聊天频道（根据您的逻辑实现）
	// 发送公会聊天消息
	// 发送公会聊天消息
	// 这里需要根据用户所在的公会找到对应的公会聊天频道进行消息分发
	// 暂时留空
}

//func (s *Server) SendPrivateChatMessage(sender *User, receiver *User, content string) {
//	// 构建私人聊天消息
//	chatMessage := &chat.ChatMessage{
//		SenderUsername:   sender.Username,
//		Content:          content,
//		ChatType:         chat.ChatMessage_ChatType_PRIVATE,
//		ReceiverUsername: receiver.Username,
//	}
//}
