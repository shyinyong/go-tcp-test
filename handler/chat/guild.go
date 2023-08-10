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
