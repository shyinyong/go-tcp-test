package gateway

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

func NewGatewayServer() *Server {
	return &Server{
		serverAddrs: make(map[protobuf.MessageType]string),
		serverConns: make(map[protobuf.MessageType]net.Conn),
	}
}

func (gs *Server) AddServers() {
	// Configure server addresses
	gs.AddServerAddr(protobuf.MessageType_LOGIN, ":8081") // Change the address
	//gs.AddServerAddr(protobuf.MessageType_GAME, ":8082")      // Change the address
	//gs.AddServerAddr(protobuf.MessageType_CHARACTER, ":8083") // Change the address
}
func (gs *Server) AddServerAddr(messageType protobuf.MessageType, address string) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.serverAddrs[messageType] = address
}

func (gs *Server) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()
	fmt.Println("Gateway server started. Listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go gs.handleGatewayConnection(conn)
	}
}

func (gs *Server) handleGatewayConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024) // Adjust the buffer size as needed
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

	// Determine the message type and forward to the appropriate server
	switch msg := clientMsg.Message.(type) {
	case *protobuf.ClientMessage_LoginRequest:
		// Handle login request
		loginReq := msg.LoginRequest
		username := loginReq.Username
		password := loginReq.Password
		fmt.Printf("received from client data： username:%s, password:%s\n", username, password)
		// Example: Assume login is successful
		successResponse := &protobuf.ServerMessage{
			Message: &protobuf.ServerMessage_LoginResponse{
				LoginResponse: &protobuf.LoginResponse{
					Success: true,
					Message: "Login successful",
				},
			},
		}

		// Serialize the response message
		responseData, err := proto.Marshal(successResponse)
		if err != nil {
			log.Println("Error marshaling response:", err)
			return
		}

		// Forward the response back to the client
		_, err = conn.Write(responseData)
		if err != nil {
			log.Println("Error writing to client connection:", err)
		}

		fmt.Println("Received login request:", msg.LoginRequest)
	case *protobuf.ClientMessage_GameActionRequest:
		// Handle game action request
		fmt.Println("Received game action request:", msg.GameActionRequest)
	case *protobuf.ClientMessage_CharacterActionRequest:
		// Handle character action request
		fmt.Println("Received character action request:", msg.CharacterActionRequest)
	case *protobuf.ClientMessage_LoginResponse:
		loginResponse := msg.LoginResponse
		// Handle login response logic, e.g., send it back to the client
		// Construct a login response message
		response := &protobuf.ServerMessage{
			Message: &protobuf.ServerMessage_LoginResponse{
				LoginResponse: &protobuf.LoginResponse{
					// Populate the fields of the response based on your logic
					Success: loginResponse.Success,
					Message: loginResponse.Message,
				},
			},
		}

		// Serialize the response message
		responseData, err := proto.Marshal(response)
		if err != nil {
			log.Println("Error marshaling response:", err)
			return
		}

		// Send the response back to the client
		_, err = conn.Write(responseData)
		if err != nil {
			log.Println("Error writing to client connection:", err)
			return
		}
	default:
		// Handle unknown message type
		fmt.Println("Received unknown message type")
	}
}

func (gs *Server) forwardMessage(messageType protobuf.MessageType, clientConn net.Conn, msg *protobuf.ClientMessage) {
	address, ok := gs.GetServerAddr(messageType)
	if !ok {
		log.Println("No server address found for message type:", messageType)
		return
	}
	log.Println("address:", address)

	serverConn, ok := gs.GetServerConn(messageType)
	if !ok {
		log.Println("No server connection found for message type:", messageType)
		return
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Println("Error marshaling message:", err)
		return
	}

	_, err = serverConn.Write(data)
	if err != nil {
		log.Println("Error writing to server connection:", err)
		return
	}

	// Receive response from server
	buffer := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := serverConn.Read(buffer)
	if err != nil {
		log.Println("Error reading from server connection:", err)
		return
	}

	// Forward the response back to the client
	_, err = clientConn.Write(buffer[:n])
	if err != nil {
		log.Println("Error writing to client connection:", err)
	}
}

func (gs *Server) GetServerAddr(messageType protobuf.MessageType) (string, bool) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	address, ok := gs.serverAddrs[messageType]
	return address, ok
}

func (gs *Server) GetServerConn(messageType protobuf.MessageType) (net.Conn, bool) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	conn, ok := gs.serverConns[messageType]
	return conn, ok
}

func (gs *Server) ConnectToServers() error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	for messageType, address := range gs.serverAddrs {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Printf("Error connecting to %s server: %v\n", messageType, err)
			return err
		}
		gs.serverConns[messageType] = conn
	}

	return nil
}
