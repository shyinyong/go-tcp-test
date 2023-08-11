package chat

import (
	"bufio"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/shyinyong/go-tcp-test/pb/chat"
	"io"
	"log"
	"net"
)

// User represents a connected user
type User struct {
	UserID           string
	Username         string
	Conn             net.Conn
	Server           *Server
	Room             *Room
	Team             *Team  // 队伍聊天
	Guild            *Guild // 公会聊天
	Writer           *bufio.Writer
	disconnectSignal chan struct{} // Signal to stop reading messages on disconnect
}

func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Conn:             conn,
		Writer:           bufio.NewWriter(conn),
		disconnectSignal: make(chan struct{}),
		Server:           server, // 设置 Server 实例
	}
	return user
}

func (u *User) disconnect() {
	log.Printf("User %s disconnected", u.Username)
	// Remove the user from the chat room
	if u.Room != nil {
		u.Room.RemoveUser(u)
	}
	close(u.disconnectSignal)
	u.Conn.Close()
}

// 监听消息处理
func (u *User) listenForMessagesWithContext(ctx context.Context) {
	reader := bufio.NewReader(u.Conn)
	for {
		select {
		case <-ctx.Done():
			return // 当连接关闭或其他情况时取消消息监听
		case <-u.disconnectSignal:
			return // 停止读取如果收到断开信号
		default:
			data, err := reader.ReadByte() // 读取一个字节
			if err != nil {
				if err != io.EOF {
					log.Println("Error reading message:", err)
				}
				break
			}

			// 根据收到的字节判断消息类型
			switch data {
			case 1: // 1 表示聊天消息
				u.handleChatMessage()
			// 还可以添加其他消息类型的处理
			default:
				log.Println("Unknown message type:", data)
			}
		}
	}
}

func (u *User) listenForMessages() {
	// 创建一个上下文，用于控制消息监听的生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 在退出方法时取消消息监听
	// 将上下文传递给 listenForMessagesWithContext
	u.listenForMessagesWithContext(ctx)
}

func (u *User) handleChatMessage() {
	// 读取消息长度
	messageBytes, err := u.readMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	// 解析 Protobuf 消息
	var chatMessage chat.ChatMessage
	err = proto.Unmarshal(messageBytes, &chatMessage)
	if err != nil {
		log.Println("Error unmarshaling chat message:", err)
		return
	}

	// 根据消息类型执行不同操作
	switch chatMessage.ChatType {
	case chat.ChatMessage_ChatType_WORLD:
		u.sendWorldChatMessage(chatMessage.Content)
	case chat.ChatMessage_ChatType_GUILD:
		u.sendGuildChatMessage(chatMessage.Content)
	case chat.ChatMessage_ChatType_PRIVATE:
		u.handlePrivateMessage(chatMessage)
	default:
		log.Println("Unknown chat type:", chatMessage.ChatType)
	}
}

func (u *User) handlePrivateMessage(chatMessage chat.ChatMessage) {
	targetUser := u.Room.FindUserByUsername(chatMessage.ReceiverUsername)
	if targetUser != nil {
		//targetUser.sendMessage(fmt.Sprintf("[Private from %s] %s", u.Username, chatMessage.Content))
		u.sendPrivateChatMessage(targetUser, chatMessage.Content)
	} else {
		u.sendMessage(fmt.Sprintf("User %s not found or not in the same room.", chatMessage.ReceiverUsername))
	}
}

// 发送普通消息
func (u *User) sendMessage(message string) {
	u.Writer.WriteString(message + "\n")
	u.Writer.Flush()
}

func (u *User) sendWorldChatMessage(content string) {
	// 构建聊天消息
	chatMessage := &chat.ChatMessage{
		SenderUsername:   u.Username,
		Content:          content,
		ChatType:         chat.ChatMessage_ChatType_WORLD,
		ReceiverUsername: "",
	}

	// 序列化消息
	messageBytes, err := proto.Marshal(chatMessage)
	if err != nil {
		log.Println("Error serializing chat message:", err)
		return
	}

	// 发送消息类型和消息内容
	u.Conn.Write([]byte{1}) // 1 表示聊天消息
	u.Conn.Write([]byte{byte(len(messageBytes) >> 24), byte(len(messageBytes) >> 16), byte(len(messageBytes) >> 8), byte(len(messageBytes))})
	u.Conn.Write(messageBytes)
}

func (u *User) sendGuildChatMessage(content string) {
	if u.Guild != nil {
		// 构建公会聊天消息
		guildChatMessage := &chat.ChatMessage{
			SenderUsername:   u.Username,
			Content:          content,
			ChatType:         chat.ChatMessage_ChatType_GUILD,
			ReceiverUsername: "",
		}

		// 序列化消息
		messageBytes, err := proto.Marshal(guildChatMessage)
		if err != nil {
			log.Println("Error serializing guild chat message:", err)
			return
		}

		// 发送消息类型和消息内容
		u.Conn.Write([]byte{1}) // 1 表示聊天消息
		u.Conn.Write([]byte{byte(len(messageBytes) >> 24), byte(len(messageBytes) >> 16), byte(len(messageBytes) >> 8), byte(len(messageBytes))})
		u.Conn.Write(messageBytes)
	}
}

func (u *User) sendPrivateChatMessage(receiver *User, content string) {
	// 构建聊天消息
	chatMessage := &chat.ChatMessage{
		SenderUsername:   u.Username,
		Content:          content,
		ChatType:         chat.ChatMessage_ChatType_PRIVATE,
		ReceiverUsername: receiver.Username,
	}

	// 序列化消息
	messageBytes, err := proto.Marshal(chatMessage)
	if err != nil {
		log.Println("Error serializing chat message:", err)
		return
	}

	// 发送消息类型和消息内容
	u.Conn.Write([]byte{1}) // 1 表示聊天消息
	u.Conn.Write([]byte{byte(len(messageBytes) >> 24), byte(len(messageBytes) >> 16), byte(len(messageBytes) >> 8), byte(len(messageBytes))})
	u.Conn.Write(messageBytes)
}

// 发送系统消息
func (u *User) sendSystemMessage(message string) {
	// 构建系统消息
	systemMessage := &chat.SystemMessage{
		Content: message,
	}

	// 序列化消息
	messageBytes, err := proto.Marshal(systemMessage)
	if err != nil {
		log.Println("Error serializing system message:", err)
		return
	}

	// 发送消息类型和消息内容
	u.Conn.Write([]byte{2}) // 2 表示系统消息
	u.Conn.Write([]byte{byte(len(messageBytes) >> 24), byte(len(messageBytes) >> 16), byte(len(messageBytes) >> 8), byte(len(messageBytes))})
	u.Conn.Write(messageBytes)
}

func (u *User) readMessage() ([]byte, error) {
	// 读取消息长度
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(u.Conn, lengthBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading message length: %w", err)
	}
	messageLength := int(lengthBytes[0])<<24 | int(lengthBytes[1])<<16 | int(lengthBytes[2])<<8 | int(lengthBytes[3])

	// 读取消息内容
	messageBytes := make([]byte, messageLength)
	_, err = io.ReadFull(u.Conn, messageBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading message content: %w", err)
	}

	return messageBytes, nil
}
