package main

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestChatServer(t *testing.T) {
	serverAddr := "localhost:8083"
	// Connect two clients to the chat server
	client1, err := net.Dial("tcp", serverAddr)
	assert.NoError(t, err)
	//defer client1.Close()
	//
	client2, err := net.Dial("tcp", serverAddr)
	//assert.NoError(t, err)
	//defer client2.Close()
	//
	//// Authenticate clients
	client1.Write([]byte("username1\n"))
	client2.Write([]byte("username2\n"))

	//
	//// Wait for authentication to complete
	//time.Sleep(time.Millisecond * 100)
	//
	//// Test user chat
	//client1.Write([]byte("Hello from client1\n"))
	//client2.Write([]byte("Hello from client2\n"))
	//
	//// Read messages from clients
	//message1 := make([]byte, 1024)
	//message2 := make([]byte, 1024)
	//
	//n1, err := client1.Read(message1)
	//assert.NoError(t, err)
	//n2, err := client2.Read(message2)
	//assert.NoError(t, err)
	//
	//assert.Contains(t, string(message1[:n1]), "Hello from client2")
	//assert.Contains(t, string(message2[:n2]), "Hello from client1")
	//
	//// Test system message
	//server.SendSystemMessage("Welcome to the chat!")
	//
	//// Read system message from clients
	//n1, err = client1.Read(message1)
	//assert.NoError(t, err)
	//n2, err = client2.Read(message2)
	//assert.NoError(t, err)
	//
	//assert.Contains(t, string(message1[:n1]), "Welcome to the chat!")
	//assert.Contains(t, string(message2[:n2]), "Welcome to the chat!")
}
