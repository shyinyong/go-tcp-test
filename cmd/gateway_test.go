package main

import (
	"fmt"
	protobuf "github.com/shyinyong/go-tcp-test/pb/message"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"testing"
)

func TestLogin(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting", err)
		return
	}
	defer conn.Close() // Ensure the connection is closed after the test finishes

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

	//// Read and handle the response from the server
	//responseData := make([]byte, 1024) // Assuming the response data can fit within this buffer size
	//n, err := conn.Read(responseData)
	//if err != nil {
	//	log.Println("Error reading response from connection:", err)
	//	return
	//}
	//
	//// Unmarshal the response data into a LoginResponse message
	//loginResponse := &protobuf.LoginResponse{}
	//err = proto.Unmarshal(responseData[:n], loginResponse)
	//if err != nil {
	//	log.Println("Error unmarshaling login response:", err)
	//	return
	//}
	//
	//// Handle the login response
	//if loginResponse.Success {
	//	fmt.Println("Login successful:", loginResponse.Message)
	//} else {
	//	fmt.Println("Login failed:", loginResponse.Message)
	//}
}
