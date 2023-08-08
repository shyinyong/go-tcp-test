package handler

import (
	protobuf "github.com/shyinyong/go-tcp-test/pb/message"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"time"
)

// HandleHeartbeat handles the client's heartbeat request and sends back a heartbeat response.
func HandleHeartbeat(conn net.Conn, clientMsg *protobuf.ClientMessage, startTime time.Time) {
	// 检查心跳请求的时间戳
	heartbeatRequest := clientMsg.GetHeartbeatRequest()
	timestamp := heartbeatRequest.GetTimestamp()

	// 计算心跳延迟
	//currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	currentTime := startTime.Unix()
	latency := currentTime - int64(timestamp)

	// 创建心跳响应消息
	heartbeatResponse := &protobuf.HeartbeatResponse{
		Latency: latency,
	}

	// 序列化心跳响应消息
	responseData, err := proto.Marshal(heartbeatResponse)
	if err != nil {
		log.Println("Error marshaling heartbeat response:", err)
		return
	}

	// 发送心跳响应消息给客户端
	_, err = conn.Write(responseData)
	if err != nil {
		log.Println("Error writing heartbeat response:", err)
		return
	}

	log.Println("Heartbeat response sent successfully")
}
