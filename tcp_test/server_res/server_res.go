package main

import (
	"bytes"
	"fmt"
	"log"
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
		fmt.Printf("Accepted connection to %v from %v\n", conn.LocalAddr(), conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Close the connection
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	// Read data from the client
	buffer := make([]byte, 10)
	var requestBuff bytes.Buffer
	for {
		n, err := conn.Read(buffer)
		fmt.Printf("Received client [%v] bytes, From %v\n", n, conn.RemoteAddr())
		if err != nil {
			fmt.Println("Error reading")
			return
		}
		requestBuff.Write(buffer[:n])
		if n < len(buffer) {
			break
		}
	}
	fmt.Printf("Received client [%s], From %v\n", requestBuff.String(), conn.RemoteAddr())

	// Write data to client
	response := "[\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0},\n{\"userId\":14939,\"event\":\"Login\",\"version\":\"13\",\"time\":\"2020-07-23 03:00:51\",\"userType\":\"guest\",\"lastLoginTime\":\"2020-07-22 23:24:49\",\"continueDays\":1,\"retentionDays\":2,\"balance\":34600381052,\"level\":362,\"vip\":1,\"ip\":\"66.188.180.40\",\"groups\":{\"Value\":1,\"Deal\":1,\"Level\":3,\"Store\":1,\"GameLayout\":1,\"Bet\":1,\"GameIndex\":3,\"AdView\":1,\"RtpSetGroup\":0},\"totalPay\":0}\n]"
	n, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Write data to client bytes:%d\n", n)
}
