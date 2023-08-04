package server

import (
	"bytes"
	"fmt"
	"github.com/shyinyong/go-tcp-test/handler"
	"io"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	var requestData bytes.Buffer
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			// Handle error
			break
		}
		requestData.Write(buf[:n])
		if n < len(buf) {
			break
		}
	}
	handler.HandleMsg(conn, requestData.Bytes())
}
