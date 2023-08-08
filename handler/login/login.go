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
	mutex sync.Mutex

	loginServerAddr string
	gameServerAddr  string
	chatServerAddr  string
}

func NewServer() *Server {
	return &Server{}
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
	//defer conn.Close()

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return
	}

	// Unmarshal the login request
	clientMsg := &protobuf.ClientMessage{}
	err = proto.Unmarshal(data[:n], clientMsg)
	if err != nil {
		log.Println("Error unmarshaling login request:", err)
		return
	}

	// Determine the message type and handle accordingly
	switch clientMsg.Type {
	case protobuf.ClientMessage_LOGIN_REQUEST:
		loginRequest := clientMsg.GetLoginRequest()
		ls.handleLoginRequest(conn, loginRequest)
	case protobuf.ClientMessage_RELOGIN_REQUEST:
		reloginRequest := clientMsg.GetReloginRequest()
		ls.handleReLoginRequest(conn, reloginRequest)
	case protobuf.ClientMessage_LOGOUT_REQUEST:
		logoutRequest := clientMsg.GetLogoutRequest()
		ls.handleLogoutRequest(conn, logoutRequest)
	default:
		log.Println("Unknown message type")
	}
}

func (ls *Server) handleLoginRequest(conn net.Conn, request *protobuf.LoginRequest) {
	// Query the database and validate the credentials
	// ... (implement your database logic here) ...
	username := request.Username
	password := request.Password
	fmt.Printf("username:%s, password:%s \n", username, password)

	// Add your login logic here
	// For example, query the database, validate credentials, etc.
	// Return the login response based on the result of your logic

	// In this example, we return a simple success response
	loginResponse := &protobuf.LoginResponse{
		Success: true,
		Message: "Login successful",
	}

	// Marshal the login response
	responseData, err := proto.Marshal(loginResponse)
	if err != nil {
		log.Println("Error marshaling login response:", err)
		return
	}

	// Send the login response back to the client
	_, err = conn.Write(responseData)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}

	// Close the connection after sending the response
	conn.Close()
}

// 重新登录
func (ls *Server) handleReLoginRequest(conn net.Conn, request *protobuf.ReLoginRequest) {
	// Add your re-login logic here
	// For example, check if the session is still valid, update session data, etc.
	// Return the re-login response based on the result of your logic

	// In this example, we return a simple success response
	reloginResponse := &protobuf.ReLoginResponse{
		Success: true,
		Message: "Re-login successful",
	}

	// Marshal the re-login response
	responseData, err := proto.Marshal(reloginResponse)
	if err != nil {
		log.Println("Error marshaling re-login response:", err)
		return
	}

	// Send the re-login response back to the client
	_, err = conn.Write(responseData)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}

// 退出
func (ls *Server) handleLogoutRequest(conn net.Conn, request *protobuf.LogoutRequest) {
	// Add your logout logic here
	// For example, invalidate the session, update user status, etc.
	// Return the logout response based on the result of your logic

	// In this example, we return a simple success response
	logoutResponse := &protobuf.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}

	// Marshal the logout response
	responseData, err := proto.Marshal(logoutResponse)
	if err != nil {
		log.Println("Error marshaling logout response:", err)
		return
	}

	// Send the logout response back to the client
	_, err = conn.Write(responseData)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}

// 登录服务器发送响应消息给网关服务器
func (ls *Server) sendResponse(conn net.Conn, responseData []byte) {
	fmt.Printf("Sending response to gateway server %v\n", responseData)
	// 这里可以添加发送响应消息给网关服务器的具体逻辑

	_, err := conn.Write(responseData)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}
