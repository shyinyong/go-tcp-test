package player

import "github.com/shyinyong/go-tcp-test/pb/cs"
import "github.com/shyinyong/go-tcp-test/msg_packet"
import "google.golang.org/protobuf/proto"

type player struct {
	cs.MsgHandlerImpl
	socket    msg_packet.WspCliConn
	accountId string
	sessionId string
	ping      int32
}

func (p *player) CMLogin(hdr *msg_packet.MsgHdr, msg *cs.CMLogin) {
	rspMsg := cs.SMLogin{}
	rspMsg.ErrCode = proto.Int32(1)
	rspMsg.ErrMsg = proto.String("invalid session_id")
}
