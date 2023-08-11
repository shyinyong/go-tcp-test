package chat

import (
	"context"
	"sync"
)

// User represents a connected user
type UserList struct {
	mutex  sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	Users  map[string]*User
}

func NewUsers() *UserList {
	ctx, cancel := context.WithCancel(context.Background())
	return &UserList{
		ctx:    ctx,
		cancel: cancel,
		Users:  make(map[string]*User),
	}
}

func (u *UserList) Len() int {
	return len(u.Users)
}
