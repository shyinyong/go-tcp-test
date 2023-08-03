package server

import (
	"go-tcp-test/config"
	"go-tcp-test/db/mysql"
	"go-tcp-test/server/socket"
)

type Server struct {
	config   *config.Config
	mysqlDB  *mysql.MySQLDB
	redisDB  *redis.RedisDB
	redisDAO *dao.RedisDAO
	handler  *socket.Handler
}

func NewServer() *Server {
	// Initialize and return a new Server instance
}

func (s *Server) Start() {
	// Start your socket server and handle incoming connections
}
