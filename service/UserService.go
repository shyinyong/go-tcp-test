package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

type UserService interface {
	Get(id string) (*User, error)
	Create(user *User) error
}

func handleUserInfoRequest(request *Request, db *sqlx.DB) (*Response, error) {
	var playerInfo PlayerInfo
	query := "SELECT uid, username, level FROM players WHERE uid = ?"
	if err := db.Get(&playerInfo, query, request.UID); err != nil {
		return nil, err
	}

	resp := &Response{
		Player: &PlayerInfo{
			UID:      playerInfo.UID,
			Username: playerInfo.Username,
			Level:    playerInfo.Level,
		},
	}

	return resp, nil
}
