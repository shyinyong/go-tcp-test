package chat

import (
	"context"
	"sync"
)

type Stat struct {
	mutex  sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	// ...
	Server *Server
}

func (st *Stat) getOnlineUserCount() int64 {
	return int64(len(st.Server.Users))
}
