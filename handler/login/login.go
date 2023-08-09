package login

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/shyinyong/go-tcp-test/config"
	protobuf "github.com/shyinyong/go-tcp-test/pb/message"
	"github.com/shyinyong/go-tcp-test/utils"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"sync"
)

type Server struct {
	mutex       sync.Mutex
	config      config.Config
	store       *sqlx.DB
	redisClient *redis.Client

	loginServerAddr string
	gameServerAddr  string
	chatServerAddr  string
}

func NewServer(cfg config.Config, store *sqlx.DB, redisClient *redis.Client) *Server {
	return &Server{
		config:      cfg,
		store:       store,
		redisClient: redisClient,
	}
}

type User struct {
	ID       string
	UserName string
}

type Session struct {
	ID   string
	User *User
}

func generateSessionID() string {
	// generate random session ID
	return "1"
}

func (ls *Server) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
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

	// Query the database for the user's information
	var storedPassword string
	err := ls.store.QueryRow("SELECT password FROM t_user WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		log.Println("Error querying database:", err)
		return
	}

	// Compare the stored password with the provided password
	var resp *protobuf.LoginResponse
	if storedPassword != password {
		log.Println("Incorrect username or password")
		// Incorrect credentials
		resp = &protobuf.LoginResponse{
			Success: false,
			Message: "Incorrect username or password",
		}
		responseData, err := proto.Marshal(resp)
		if err != nil {
			log.Println("Error marshaling login response:", err)
			return
		}
		ls.sendResponse(conn, responseData)
		return
	}

	// Generate a session ID and save it in Redis
	// Implement this function to generate a unique session ID
	sessionID, _ := utils.GenerateSessionID()
	err = ls.saveSessionToRedis(sessionID, request.Username) // Save the session ID along with the username
	if err != nil {
		log.Println("Error saving session to Redis:", err)
		return
	}

	resp = &protobuf.LoginResponse{
		Success: true,
		Message: "Login successful",
	}

	// Marshal the login response
	responseData, err := proto.Marshal(resp)
	if err != nil {
		log.Println("Error marshaling login response:", err)
		return
	}

	// Send the login response back to the client
	ls.sendResponse(conn, responseData)
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
	// 这里可以添加发送响应消息给网关服务器的具体逻辑
	_, err := conn.Write(responseData)
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}

// Save the session ID and username in Redis
func (ls *Server) saveSessionToRedis(sessionID, username string) error {
	ctx := context.Background()
	err := ls.redisClient.Set(ctx, sessionID, username, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Example function to validate session using Redis
func (ls *Server) validateSession(sessionID string) (bool, error) {
	ctx := context.Background()
	username, err := ls.redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Session not found
		}
		return false, err // Error occurred
	}

	// Session found
	_ = username // You can use the username for further processing if needed
	return true, nil
}
