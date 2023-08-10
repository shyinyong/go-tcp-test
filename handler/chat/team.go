package chat

import (
	"fmt"
	"sync"
)

type Team struct {
	ID    int
	Name  string
	Users map[*User]bool
	mu    sync.Mutex
}

func NewTeam(id int, name string) *Team {
	return &Team{
		ID:    id,
		Name:  name,
		Users: make(map[*User]bool),
	}
}

func (t *Team) AddUser(user *User) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Users[user] = true
	user.Team = t
}

func (t *Team) RemoveUser(user *User) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.Users, user)
}

func (t *Team) Broadcast(sender *User, message string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for user := range t.Users {
		if user != sender {
			user.sendMessage(fmt.Sprintf("[Team %s] %s: %s", t.Name, sender.Username, message))
		}
	}
}
