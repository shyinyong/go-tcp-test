package tcp_test

//
//import (
//	"bytes"
//	"client_app/buffer"
//	"fmt"
//	"io"
//	"net"
//	"time"
//)
//
//func login(conn net.Conn) {
//	loginData := buffer.GenLoginData()
//	_, err := conn.Write(loginData.Bytes())
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("用户登录成功......")
//}
//
//func recv(c chan int) {
//	receiveData := <-c
//	fmt.Printf("receive:%d\n", receiveData)
//}
//
//func send(c chan int, n int) {
//	c <- n
//}
//
//func spinner(delay time.Duration) {
//	for {
//		for _, r := range `-|/` {
//			fmt.Printf("\r%c", r)
//			time.Sleep(delay)
//		}
//	}
//}
//
//func sendRequest(conn net.Conn, cmdID int) {
//	data := buffer.GenCMDIDData(cmdID)
//	conn.Write(data.Bytes())
//}
//
//func onMessageRecived(conn *net.TCPConn) {
//	for {
//		buf := make([]byte, 1024)
//		cnt, err := conn.Read(buf)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		fmt.Println("服务器返回信息：", string(buf[:cnt]))
//		time.Sleep(time.Second * 1)
//	}
//}
//
//const bufferLen = 1024 // 缓冲区大小
//const recvLen = 1024   // 期望接收的数据长度
//
//func readData(conn net.Conn) ([]byte, error) {
//	buffer := make([]byte, bufferLen)
//	data := make([]byte, 0, recvLen)
//	dataBuffer := bytes.Buffer{}
//
//	// 设置读取超时时间为5秒
//	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
//
//	for {
//		fmt.Printf("len(data):%d\n", len(data))
//		n, err := conn.Read(buffer)
//		if err != nil {
//			return nil, err
//		}
//		if n > 0 {
//			dataBuffer.Write(buffer[:n])
//		}
//
//		// 如果连接已经关闭，则返回读取数据
//		if err == io.EOF {
//			return dataBuffer.Bytes(), nil
//		}
//
//		// 如果读取超时，判断是否读取完成，如果已经读取完成则返回数据，否则返回错误
//		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
//			if n == 0 {
//				return nil, fmt.Errorf("read data timeout")
//			} else {
//				continue
//			}
//		}
//
//		// 将读取到的数据添加到 data 中
//		//data = append(data, buffer[:n]...)
//		// 如果读取到的数据已经超过了期望的长度，忽略多余的部分
//		//if len(data) > recvLen {
//		//	data = data[:recvLen]
//		//}
//		// 如果读取到的数据已经等于期望的长度，说明读取完成，退出循环
//		//if len(data) == recvLen {
//		//	break
//		//}
//	}
//	return dataBuffer.Bytes(), nil
//}
//
//const (
//	CTSActive State = iota + 1
//	CTSRevoked
//	CTSUnknown
//)
//
//type State int32
//
//func main() {
//	addr, _ := net.ResolveTCPAddr("tcp", "192.168.1.207:7235")
//	conn, err := net.DialTCP("tcp", nil, addr)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer conn.Close()
//	fmt.Println("连接服务器成功......")
//
//	// 发送登录数据包
//	login(conn)
//	loginBuffData, err := readData(conn)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("login data:", string(loginBuffData))
//
//	// Next
//	cmdID := 1656
//	sendRequest(conn, cmdID)
//	fmt.Printf("cmdID:%d\n", cmdID)
//	buffData, err := readData(conn)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("cmd data:", string(buffData))
//}
