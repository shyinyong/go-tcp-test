package chat

import (
	"fmt"
	"sync"
)

var defaultChatRoom *ChatRoom

type ChatRoom struct {
	messages    []string
	mu          sync.Mutex
	Name        string
	Users       map[*User]bool
	Teams       map[*Team]bool
	broadcast   chan string
	systemMsg   string
	systemMsgMu sync.Mutex
}

func NewChatRoom(name string) *ChatRoom {
	room := &ChatRoom{
		Name:      name,
		Users:     make(map[*User]bool),
		broadcast: make(chan string),
	}
	go room.start()
	return room
}

// FindUserByUsername finds a user in the same room by username
func (r *ChatRoom) FindUserByUsername(username string) *User {
	r.mu.Lock()
	defer r.mu.Unlock()

	for user := range r.Users {
		if user.Username == username {
			return user
		}
	}

	return nil
}

// GetOnlineUsers returns a list of usernames of online users in the room
func (r *ChatRoom) GetOnlineUsers() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	var onlineUsers []string
	for user := range r.Users {
		onlineUsers = append(onlineUsers, user.Username)
	}

	return onlineUsers
}

func (r *ChatRoom) RemoveUser(user *User) {
	delete(r.Users, user)
}

func (r *ChatRoom) Broadcast(sender *User, message string) {
	// Send the message to the broadcast channel
	r.broadcast <- fmt.Sprintf("[%s] %s: %s", r.Name, sender.Username, message)
}

func (r *ChatRoom) start() {
	for {
		select {
		case message := <-r.broadcast:
			// Broadcast the message to all users in the room
			for user := range r.Users {
				user.Conn.Write([]byte(message + "\n"))
			}
		}
	}
}

func GetDefaultChatRoom() *ChatRoom {
	if defaultChatRoom == nil {
		defaultChatRoom = NewChatRoom("DefaultRoom")
	}
	return defaultChatRoom
}

func (r *ChatRoom) AddUser(user *User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Users[user] = true
	user.Room = r

	//r.Users[user] = true
}

// BroadcastSystemMessage broadcasts a system message to all users in the room
func (r *ChatRoom) BroadcastSystemMessage(message string) {
	r.systemMsgMu.Lock()
	r.systemMsg = message
	r.systemMsgMu.Unlock()

	r.broadcast <- message
}

// GetUserList returns the list of usernames of users in the room
func (r *ChatRoom) GetUserList() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	var userList []string
	for user := range r.Users {
		userList = append(userList, user.Username)
	}
	return userList
}

func (r *ChatRoom) FindTeamByID(teamID int) *Team {
	r.mu.Lock()
	defer r.mu.Unlock()

	for team := range r.Teams {
		if team.ID == teamID {
			return team
		}
	}

	return nil
}
