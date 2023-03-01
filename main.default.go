package main

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// 	"net"
// 	"os"
// 	"time"
// )

// const (
// 	connHost = "192.168.1.207"
// 	connPort = "7235"
// 	connType = "tcp"
// )

// type Pkt struct {
// 	CmdId uint16
// }

// func makeLoginPktData() *bytes.Buffer {
// 	//account := "m-127054pkgk"
// 	//password := "74516b5c98fe"

// 	account := "m-127824qspu"
// 	password := "2371588288ba"

// 	// account := "m-127750aqya"
// 	// password := "ede2fc7d1c69"

// 	pkt := new(bytes.Buffer)
// 	binary.Write(pkt, binary.LittleEndian, uint16(1000))
// 	binary.Write(pkt, binary.LittleEndian, uint16(1700))
// 	binary.Write(pkt, binary.LittleEndian, uint8(len(account)))
// 	binary.Write(pkt, binary.LittleEndian, []byte(account))
// 	binary.Write(pkt, binary.LittleEndian, uint8(len(password)))
// 	binary.Write(pkt, binary.LittleEndian, []byte(password))
// 	binary.Write(pkt, binary.LittleEndian, uint16(19))
// 	binary.Write(pkt, binary.LittleEndian, int32(1))

// 	data := new(bytes.Buffer)
// 	binary.Write(data, binary.LittleEndian, uint16(pkt.Len()))
// 	binary.Write(data, binary.LittleEndian, pkt.Bytes())

// 	return data
// 	// fmt.Print(data.Bytes())
// 	// return
// }

// func makeClientPktData() *bytes.Buffer {
// 	cmdId := 1659
// 	var param1 uint32 = 3

// 	// pkt
// 	pkt := new(bytes.Buffer)
// 	binary.Write(pkt, binary.LittleEndian, uint16(cmdId))
// 	binary.Write(pkt, binary.LittleEndian, uint32(param1))
// 	// data
// 	data := new(bytes.Buffer)
// 	binary.Write(data, binary.LittleEndian, uint16(pkt.Len()))
// 	binary.Write(data, binary.LittleEndian, pkt.Bytes())

// 	return data
// }

// func apiClient() {
// 	fmt.Println("Connecting to " + connType + " server " + connHost + ":" + connPort)
// 	conn, err := net.Dial(connType, connHost+":"+connPort)
// 	if err != nil {
// 		fmt.Println("Error connecting:", err.Error())
// 		os.Exit(1)
// 	}
// 	defer conn.Close() //关闭

// 	login_data := makeLoginPktData()
// 	_, err = conn.Write(login_data.Bytes())
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// login response
// 	loginRecvData := make([]byte, 5000)
// 	n, err := conn.Read(loginRecvData) //读取数据
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	recvStr := string(loginRecvData[:n])
// 	fmt.Printf("Response data:\n %s \n", recvStr)

// 	fmt.Println(" next ------------------------ next ")
// 	time.Sleep(1 * time.Second)

// 	// Custom api

// 	client_data := makeClientPktData()
// 	_, err = conn.Write(client_data.Bytes())
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// fmt.Printf("Send data success: %b", client_data.Bytes())
// 	// for {
// 	// 	time.Sleep(1 * time.Second)
// 	// }

// 	// response
// 	// recvData := make([]byte, 5000)
// 	// _, err = conn.Read(recvData) //读取数据
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return
// 	// }
// 	// recvStr = string(recvData[:4])
// 	// fmt.Printf("Response data: %s", recvStr)

// }

// func main() {
// 	// for i := 1; i <= 20; i++ {
// 	// 	fmt.Println("round:", i)
// 	// 	apiClient()
// 	// 	time.Sleep(1 * time.Second)
// 	// }

// 	round := 0
// 	for {
// 		round++
// 		fmt.Println("round:", round)
// 		apiClient()
// 		// time.Sleep(1 * time.Second)

// 		if round >= 1 {
// 			fmt.Println("finished")
// 			break
// 		}
// 	}
// }
