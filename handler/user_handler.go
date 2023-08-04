package handler

import (
	"github.com/shyinyong/go-tcp-test/db/mysql"
	"github.com/shyinyong/go-tcp-test/pb/request"
	pb "github.com/shyinyong/go-tcp-test/pb/user"
)

db *
func UserInfo(msg *request.Request) {

}

func GetUserInfo(msg *request.Request)  {
	db := mysql.GetDB()
	var playerInfo pb.PlayerInfo
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