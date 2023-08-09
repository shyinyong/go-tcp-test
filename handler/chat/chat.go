package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	mutex sync.Mutex
}

func NewServer() *Server {
	return &Server{}
}

func (cs *Server) Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	fmt.Println("Chat server started. Listening on", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go cs.handleChatConnection(conn)
	}
}

func (cs *Server) handleChatConnection(conn net.Conn) {
	data := make([]byte, 1024)
	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Println("Error reading from connection:", err)
			conn.Close()
			return
		}

		message := string(data[:n])
		fmt.Println("Received message:", message)

		// Implement your chat logic here

		// Echo the received message back to the client
		_, err = conn.Write([]byte("You said: " + message))
		if err != nil {
			log.Println("Error writing to connection:", err)
			conn.Close()
			return
		}
	}
}
