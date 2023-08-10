package chat

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/shyinyong/go-tcp-test/pb/chat"
	"io"
	"log"
	"net"
)

// User represents a connected user
type User struct {
	Username              string
	Conn                  net.Conn
	Server                *Server
	Room                  *ChatRoom
	Team                  *Team  // 队伍聊天
	Guild                 *Guild // 公会聊天
	PrivateMessageChannel chan *chat.ChatMessage
	Writer                *bufio.Writer
	disconnectSignal      chan struct{} // Signal to stop reading messages on disconnect
}

func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Conn:                  conn,
		Writer:                bufio.NewWriter(conn),
		disconnectSignal:      make(chan struct{}),
		Server:                server, // 设置 Server 实例
		PrivateMessageChannel: make(chan *chat.ChatMessage),
	}
	go user.handlePrivateMessages() // 启动私聊消息处理协程
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

func (u *User) listenForMessages() {
	reader := bufio.NewReader(u.Conn)
	for {
		select {
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

func (u *User) handleChatMessage() {
	// 读取消息长度
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(u.Conn, lengthBytes)
	if err != nil {
		log.Println("Error reading message length:", err)
		return
	}
	messageLength := int(lengthBytes[0])<<24 | int(lengthBytes[1])<<16 | int(lengthBytes[2])<<8 | int(lengthBytes[3])

	// 读取消息内容
	messageBytes := make([]byte, messageLength)
	_, err = io.ReadFull(u.Conn, messageBytes)
	if err != nil {
		log.Println("Error reading message content:", err)
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
		// 处理公会聊天
	case chat.ChatMessage_ChatType_PRIVATE:
		u.handlePrivateMessage(chatMessage)
	default:
		log.Println("Unknown chat type:", chatMessage.ChatType)
	}
}

func (u *User) handlePrivateMessage(chatMessage chat.ChatMessage) {
	targetUser := u.Room.FindUserByUsername(chatMessage.ReceiverUsername)
	if targetUser != nil {
		targetUser.sendMessage(fmt.Sprintf("[Private from %s] %s", u.Username, chatMessage.Content))
	} else {
		u.sendMessage(fmt.Sprintf("User %s not found or not in the same room.", chatMessage.ReceiverUsername))
	}
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

// 发送聊天消息
func (u *User) sendChatMessage(content string, chatType chat.ChatMessage_ChatType, receiverUsername string) {
	// 构建聊天消息
	message := &chat.ChatMessage{
		SenderUsername:   u.Username,
		Content:          content,
		ChatType:         chatType,
		ReceiverUsername: receiverUsername,
	}

	// 序列化消息
	messageBytes, err := proto.Marshal(message)
	if err != nil {
		log.Println("Error serializing chat message:", err)
		return
	}

	// 发送消息类型和消息内容
	u.Conn.Write([]byte{1}) // 1 表示聊天消息
	u.Conn.Write([]byte{byte(len(messageBytes) >> 24), byte(len(messageBytes) >> 16), byte(len(messageBytes) >> 8), byte(len(messageBytes))})
	u.Conn.Write(messageBytes)
}

func (u *User) handleMessage(message *chat.ChatMessage) {
	// 处理私聊消息
	if u.Room == u.Server.worldRoom {
		u.Server.SendWorldChatMessage(u, message.GetContent())
	} else {
		// 这里可以根据消息类型执行不同的操作
		// 比如处理私聊消息
		receiver := u.Room.FindUserByUsername(message.GetReceiverUsername())
		if receiver != nil {
			u.SendPrivateMessage(receiver, message.GetContent())
		} else {
			u.sendMessage("User not found or not in the same room.")
		}
	}
}

// SendPrivateMessage 发送私聊消息
func (u *User) SendPrivateMessage(receiver *User, content string) {
	// 创建私聊消息
	privateMsg := &chat.ChatMessage{
		SenderUsername:   u.Username,
		ReceiverUsername: receiver.Username,
		Content:          content,
	}
	// 将私聊消息发送给目标用户
	receiver.PrivateMessageChannel <- privateMsg
}

func (u *User) sendMessage(message string) {
	u.Writer.WriteString(message + "\n")
	u.Writer.Flush()
}

func (u *User) sendTeamMessage(message string) {
	if u.Team != nil {
		u.Team.Broadcast(u, message)
	}
}

// handlePrivateMessages 处理私聊消息
func (u *User) handlePrivateMessages() {
	for {
		select {
		case <-u.disconnectSignal:
			return
		case privateMsg := <-u.PrivateMessageChannel:
			// 处理接收到的私聊消息
			u.handleReceivedPrivateMessage(privateMsg)
		}
	}
}

func (u *User) handleReceivedPrivateMessage(privateMsg *chat.ChatMessage) {
	// 在这里处理接收到的私聊消息，例如将消息发送给用户、记录日志等
	fmt.Printf("handleReceivedPrivateMessage:%s", privateMsg.Content)
}
