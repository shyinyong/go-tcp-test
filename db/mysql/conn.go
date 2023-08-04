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

var DB *sqlx.DB

func Init(config *config.Config) {
	var err error
	DB, err = sqlx.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer DB.Close()
}

func GetDB() *sqlx.DB {
	return DB
}

//func (m *MySQLDB) Close() error {
//	return m.db.Close()
//}
