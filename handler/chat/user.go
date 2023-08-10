package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// User represents a connected user
type User struct {
	Username         string
	Conn             net.Conn
	Room             *ChatRoom
	Writer           *bufio.Writer
	disconnectSignal chan struct{} // Signal to stop reading messages on disconnect
}

func NewUser(conn net.Conn) *User {
	return &User{
		Conn:             conn,
		Writer:           bufio.NewWriter(conn),
		disconnectSignal: make(chan struct{}),
	}
}

func (u *User) listenForMessages() {
	reader := bufio.NewReader(u.Conn)
	for {
		select {
		case <-u.disconnectSignal:
			return // Stop reading if disconnect signal is received
		default:
			message, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Println("Error reading message:", err)
				}
				break
			}

			// Handle the received message
			u.handleMessage(message)
		}
	}
}

func (u *User) handleMessage(message string) {
	// Check if the message is a private message
	if strings.HasPrefix(message, "/msg ") {
		u.sendPrivateMessage(message)
		return
	}

	// Forward the message to the chat room for distribution
	if u.Room != nil {
		u.Room.Broadcast(u, message)
	}
}

func (u *User) sendPrivateMessage(message string) {
	// Parse the private message format: "/msg <username> <message>"
	parts := strings.SplitN(message, " ", 3)
	if len(parts) != 3 {
		u.sendMessage("Invalid private message format. Usage: /msg <username> <message>")
		return
	}

	targetUsername := parts[1]
	targetMessage := parts[2]

	// Find the target user by username
	targetUser := u.Room.FindUserByUsername(targetUsername)
	if targetUser == nil {
		u.sendMessage(fmt.Sprintf("User %s not found or not in the same room.", targetUsername))
		return
	}

	// Send the private message to the target user
	targetUser.sendMessage(fmt.Sprintf("[Private from %s] %s", u.Username, targetMessage))
}

func (u *User) sendMessage(message string) {
	u.Writer.WriteString(message + "\n")
	u.Writer.Flush()
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
