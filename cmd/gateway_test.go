package main

import (
	"fmt"
	protobuf "github.com/shyinyong/go-tcp-test/pb/message"
	"google.golang.org/protobuf/proto"
	"log"
	"testing"
)

func TestLogin(t *testing.T) {
	conn := newTestServer(t)

	// Create a login request message
	loginRequest := &protobuf.LoginRequest{
		Username: "apple",
		Password: "123456",
	}
	clientMsg := &protobuf.ClientMessage{
		Type: protobuf.ClientMessage_LOGIN_REQUEST,
		Message: &protobuf.ClientMessage_LoginRequest{
			LoginRequest: loginRequest,
		},
	}

	// Serialize the message
	data, err := proto.Marshal(clientMsg)
	if err != nil {
		log.Println("Error marshaling message:", err)
		return
	}

	// Send the message to the gateway server
	_, err = conn.Write(data)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}

	fmt.Println("Login request sent to gateway server.")
}
