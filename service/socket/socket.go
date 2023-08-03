package socket

import (
	"your_project/internal/database/mysql"
	"your_project/internal/database/redis/dao"
)

type Handler struct {
	mysqlDB  *mysql.MySQLDB
	redisDAO *dao.RedisDAO
}

func NewHandler(mysqlDB *mysql.MySQLDB, redisDAO *dao.RedisDAO) *Handler {
	return &Handler{mysqlDB: mysqlDB, redisDAO: redisDAO}
}
