package tcp_test

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	message := "Hello from the client!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := make([]byte, 1024)
	var response bytes.Buffer
	for {
		n, err := conn.Read(buffer)
		fmt.Printf("Received server bytes:%d, and buffer len:%d \n", n, len(buffer))
		if err != nil {
			fmt.Println("Error reading")
			return
		}
		response.Write(buffer[:n])
		if n < len(buffer) {
			break
		}
	}
	fmt.Printf("Received server data:%s\n", response.String())
}

func checkErr(err error) {
	if err != nil {
		log.Fatal()
	}
}
