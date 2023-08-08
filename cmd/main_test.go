package main

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func newTestServer(t *testing.T) (conn net.Conn) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting", err)
	}

	return conn
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
