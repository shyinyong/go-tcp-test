package service

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jmoiron/sqlx"
	_ "google.golang.org/protobuf/proto"
	"time"
)

type PlayerInfo struct {
	UID      int32  `db:"uid"`
	Username string `db:"username"`
	Level    int32  `db:"level"`
}

type User struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Occupation string    `json:"occupation,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserService struct {
}

func getUserInfo() {

}
