package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/shyinyong/go-tcp-test/config"
)

var Db *sqlx.DB

func NewDB(config *config.Config) *sqlx.DB {
	Db, err := sqlx.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err)
	}
	//defer Db.Close()
	return Db
}

func Close() error {
	return Db.Close()
}
