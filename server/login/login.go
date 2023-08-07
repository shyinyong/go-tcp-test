package login

import (
	"fmt"
	protobuf "github.com/shyinyong/go-tcp-test/pb/message"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"sync"
)

type Server struct {
	mutex       sync.Mutex
	serverAddrs map[protobuf.MessageType]string
	serverConns map[protobuf.MessageType]net.Conn
}

func NewLoginServer() *Server {
	return &Server{
		serverAddrs: make(map[protobuf.MessageType]string),
		serverConns: make(map[protobuf.MessageType]net.Conn),
	}
}

func (ls *Server) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Login server started. Listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go ls.handleLoginConnection(conn)
	}
}

func (ls *Server) handleLoginConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return
	}

	clientMsg := &protobuf.ClientMessage{}
	err = proto.Unmarshal(buffer[:n], clientMsg)
	if err != nil {
		log.Println("Error unmarshaling message:", err)
		return
	}

	// Check if the received message is a LoginRequest
	loginReq, ok := clientMsg.Message.(*protobuf.ClientMessage_LoginRequest)
	if !ok {
		log.Println("Invalid message type. Expected LoginRequest.")
		return
	}

	// Extract username and password from the login request
	username := loginReq.LoginRequest.Username
	password := loginReq.LoginRequest.Password
	fmt.Println("username:%s, password:%s", username, password)

	// Perform authentication and database query here
	// ...

	// Simulate a successful login
	loginResponse := &protobuf.ServerMessage{
		Message: &protobuf.ServerMessage_LoginResponse{
			LoginResponse: &protobuf.LoginResponse{
				Success: true,
				Message: "Login successful",
			},
		},
	}

	// Serialize the response message
	responseData, err := proto.Marshal(loginResponse)
	if err != nil {
		log.Println("Error marshaling response:", err)
		return
	}

	// Send the response back to the client
	_, err = conn.Write(responseData)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}
