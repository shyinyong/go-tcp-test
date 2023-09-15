package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var err error

func InitDB() {
	// argv config *config.Config
	// Db, err := sqlx.Open(config.DBDriver, config.DBSource)
	db, err = sqlx.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/go_tcp_test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		glog.Fatal(err)
	}
}

func GetDB() *sqlx.DB {
	return db
}

func Close() error {
	return db.Close()
}
