package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started, waiting for connections...")

	// Accept and handle incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Read data from the client
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Print(err)
		return
	}

	data := buffer[:n]
	fmt.Printf("Received data:%s\n", data)

	response := "[\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0}\n]"
	rLen := len(response)
	n, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Server data len:%d, and send len:%d\n", rLen, n)

	// Close the connection
	conn.Close()
}
