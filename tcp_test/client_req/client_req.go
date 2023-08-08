package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/shyinyong/go-tcp-test/tcp_test/msg"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.207:7235")
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

	// step1, login request
	message := msg.GenLoginData().Bytes()
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	// step1, login response
	buffer := make([]byte, 1024)
	var response bytes.Buffer
	for {
		n, err := conn.Read(buffer)
		fmt.Printf("Received server bytes:%d, and msg len:%d \n", n, len(buffer))
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

	time.Sleep(5)

	// step2, custom msg_id request
	cmdID := msg.GenCMDIDData(1656)
	writeResponse(conn, cmdID.Bytes())
	fmt.Printf("cmdID:%d\n", cmdID)
	buffData, err := readRequest(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("cmd data:", string(buffData))

}

func readRequest(conn net.Conn) ([]byte, error) {
	// 读取消息头中的长度字段
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lenBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf)

	// 读取消息体
	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil

	//var requestData = bytes.Buffer{}
	//buf := make([]byte, 256)
	//for {
	//	n, err := conn.Read(buf)
	//	if err != nil {
	//		if err != io.EOF {
	//			fmt.Println("read error:", err)
	//		}
	//		// Handle error
	//		break
	//	}
	//	requestData.Write(buf[:n])
	//	if n < len(buf) {
	//		break
	//	}
	//}
	//return requestData.Bytes(), nil
}

func writeResponse(conn net.Conn, data []byte) error {
	// 构造消息头
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))

	// 发送消息头和消息体
	_, err := conn.Write(append(header, data...))
	return err
}

func checkErr(err error) {
	if err != nil {
		log.Fatal()
	}
}
