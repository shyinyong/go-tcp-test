package main

import (
	"bytes"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting", err)
		return
	}

	defer conn.Close()

	message := "Hello from the client!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Received response from the server
	buffer := make([]byte, 1024)
	var response bytes.Buffer
	for {
		n, err := conn.Read(buffer)
		fmt.Printf("Received response len:%d, and buffer len:%d \n", n, len(buffer))
		if err != nil {
			fmt.Println("Error reading")
			return
		}

		response.Write(buffer[:n])
		if n < len(buffer) {
			break
		}
	}
	fmt.Printf("Received response:%s\n", response.String())
}
