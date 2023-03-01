package main

import (
	"bufio"
	"client_app/buffer"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func login(conn net.Conn) {
	recv, send := make(chan string), make(chan string)
	go func(s chan string) {
		loginData := buffer.GenLoginData()
		_, err := conn.Write(loginData.Bytes())
		if err != nil {
			panic(err)
		}
		fmt.Println("用户登录成功......")
	}(send)
	go func(r chan string) {
		buf := make([]byte, 1024)
		cnt, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		r <- fmt.Sprintf("recv: %v", string(buf[:cnt]))
	}(recv)
	select {
	case accept := <-recv:
		fmt.Println("accept")
		log.Println(accept)
	case to := <-send:
		fmt.Println("send to ")
		log.Println(to)
	}

}

func recv(c chan int) {
	receiveData := <-c
	fmt.Printf("receive:%d\n", receiveData)
}

func send(c chan int, n int) {
	c <- n
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func sendRequest(conn net.Conn, cmdID uint16) {
	data := buffer.GenCMDIDData(cmdID)
	conn.Write(data.Bytes())

	//读取服务器返回的信息
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	fmt.Println("等待服务器返回信息 len=", n)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("服务器返回信息：", string(buf[:n]))
}

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "192.168.1.207:7235")
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("连接服务器成功......")

	// 发送登录数据包
	login(conn)
	fmt.Println("用户登录成功......")

	fmt.Println("Please input the cmdID: ")
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error input:", err)
			os.Exit(1)
		}
		line = strings.Trim(line, " \r\n")
		if len(line) == 0 {
			fmt.Println("发送消息id不能为空")
			continue
		}
		cmdID := buffer.StrConvToUint16(line)
		if cmdID == 0 {
			break
		}

		sendRequest(conn, cmdID)
		fmt.Println("Please input the cmdID: ")
	}
}
