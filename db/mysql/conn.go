package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/shyinyong/go-tcp-test/config"
)

//type MySQLDB struct {
//	db *sql.DB
//}

//func NewMySQLDB(config *config.Config) (*MySQLDB, error) {
//	db, err := sqlx.Open(config.DBDriver, config.DBSource)
//	if err != nil {
//		return nil, err
//	}
//	defer db.Close()
//
//	return &MySQLDB{db: db}, nil
//}

//
//func (m *MySQLDB) Close() error {
//	return m.db.Close()
//}

var Db *sqlx.DB

func NewDB(config *config.Config) *sqlx.DB {
	Db, err := sqlx.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err)
	}
	//defer Db.Close()
	return Db
}
