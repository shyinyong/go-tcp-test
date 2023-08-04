package server

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"sync"
)

type Server struct {
	mutex sync.Mutex
}

func Start() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer listener.Close()
	log.Info().Msg("Server started, waiting for connections...")

	var wg sync.WaitGroup
	server := &Server{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal().Err(err)
			continue
		}
		fmt.Printf("Accepted connection to %v from %v\n", conn.LocalAddr(), conn.RemoteAddr())

		//go server.handleRequest(conn)
		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			defer conn.Close()
			server.handleConnection(conn)
		}(conn)
	}

	wg.Wait()
}
