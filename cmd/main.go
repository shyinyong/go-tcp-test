package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shyinyong/go-tcp-test/config"
	"github.com/shyinyong/go-tcp-test/db/mysql"
	"github.com/shyinyong/go-tcp-test/server"
	"os"
)

func main() {
	// Config env initialize
	cfg, err := config.LoadConfig(".")
	checkErr(err)
	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Initialize database connections
	mysql.Init(&cfg)

	// Start the server
	server.Start()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal()
	}
}
