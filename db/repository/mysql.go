package repository

import (
	"github.com/shyinyong/go-tcp-test/db/mysql"
)

type MySQLRepository struct {
	db *mysql.MySQLDB
}

func NewMySQLRepository(db *mysql.MySQLDB) *MySQLRepository {
	return &MySQLRepository{db: db}
}
