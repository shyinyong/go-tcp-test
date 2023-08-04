package service

import (
	"github.com/shyinyong/go-tcp-test/config"
	"github.com/shyinyong/go-tcp-test/db/mysql"
	"github.com/shyinyong/go-tcp-test/db/redis"
	"github.com/shyinyong/go-tcp-test/service/socket"
)

type Server struct {
	config  *config.Config
	mysqlDB *mysql.MySQLDB
	redisDB *redis.DB
	handler *socket.Handler
}

func NewServer() *Server {
	// Initialize and return a new Server instance
	return &Server{}
}

func (s *Server) Start() {
	// Start your socket server and handle incoming connections
}
