package chat

import (
	"fmt"
	"sync"
)

type Guild struct {
	ID    int
	Name  string
	Users map[*User]bool
	mu    sync.Mutex
}

func NewGuild(id int, name string) *Guild {
	return &Guild{
		ID:    id,
		Name:  name,
		Users: make(map[*User]bool),
	}
}

func (g *Guild) AddUser(user *User) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Users[user] = true
	user.Guild = g
}

func (g *Guild) RemoveUser(user *User) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.Users, user)
}

func (g *Guild) Broadcast(sender *User, message string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for user := range g.Users {
		if user != sender {
			user.sendMessage(fmt.Sprintf("[Guild %s] %s: %s", g.Name, sender.Username, message))
		}
	}
}

// 如果您的设计是每个公会都有一个聊天频道，而不是每个用户都有一个公会聊天频道属性，那么可以从服务器层级来管理公会聊天频道，
// 并在服务器的 rooms 中维护公会聊天频道。然后，在发送公会聊天消息时，直接使用公会聊天频道来进行广播。
