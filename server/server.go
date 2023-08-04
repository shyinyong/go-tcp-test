package server

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"net"
)

func Start() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer listener.Close()
	log.Info().Msg("Server started, waiting for connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal().Err(err)
		}

		fmt.Printf("Accepted connection to %v from %v\n", conn.LocalAddr(), conn.RemoteAddr())
		go handleConn(conn)
	}
}