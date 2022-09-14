package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	connHost = "192.168.1.207"
	connPort = "7235"
	connType = "tcp"
)

type Pkt struct {
	CmdId uint16
}

func makeLoginPktData() *bytes.Buffer {
	account := "m-126108hgnn"
	password := "a58870025879"

	pkt := new(bytes.Buffer)
	binary.Write(pkt, binary.LittleEndian, uint16(1000))
	binary.Write(pkt, binary.LittleEndian, uint16(1700))
	binary.Write(pkt, binary.LittleEndian, uint8(len(account)))
	binary.Write(pkt, binary.LittleEndian, []byte(account))
	binary.Write(pkt, binary.LittleEndian, uint8(len(password)))
	binary.Write(pkt, binary.LittleEndian, []byte(password))
	binary.Write(pkt, binary.LittleEndian, uint16(19))
	binary.Write(pkt, binary.LittleEndian, int32(1))

	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, uint16(pkt.Len()))
	binary.Write(data, binary.LittleEndian, pkt.Bytes())

	return data
	// fmt.Print(data.Bytes())
	// return
}

func makeClientPktData() *bytes.Buffer {
	cmdId := 1612
	var param1 uint32 = 1

	// pkt
	pkt := new(bytes.Buffer)
	binary.Write(pkt, binary.LittleEndian, uint16(cmdId))
	binary.Write(pkt, binary.LittleEndian, uint32(param1))
	// data
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, uint16(pkt.Len()))
	binary.Write(data, binary.LittleEndian, pkt.Bytes())

	return data
}

func main() {
	fmt.Println("Connecting to " + connType + " server " + connHost + ":" + connPort)
	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close() //关闭

	login_data := makeLoginPktData()
	_, err = conn.Write(login_data.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}

	// login response
	loginRecvData := make([]byte, 5000)
	n, err := conn.Read(loginRecvData) //读取数据
	if err != nil {
		fmt.Println(err)
		return
	}
	recvStr := string(loginRecvData[:n])
	fmt.Printf("Response data: %s \n", recvStr)

	fmt.Println(" ------------------------ ")

	// custom api
	time.Sleep(1 * time.Second)

	client_data := makeClientPktData()
	_, err = conn.Write(client_data.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("Send data success: %b", client_data.Bytes())
	// for {
	// 	time.Sleep(1 * time.Second)
	// }

	// response
	recvData := make([]byte, 5000)
	_, err = conn.Read(recvData) //读取数据
	if err != nil {
		fmt.Println(err)
		return
	}
	recvStr = string(recvData[:4])
	fmt.Printf("Response data: %s", recvStr)

}
