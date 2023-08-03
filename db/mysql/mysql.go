package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shyinyong/go-tcp-test/config"
)

type MySQLDB struct {
	db *sql.DB
}

func NewMySQLDB(config *config.Config) (*MySQLDB, error) {
	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return nil, err
	}

	return &MySQLDB{db: db}, nil
}

func (m *MySQLDB) Close() error {
	return m.db.Close()
}
