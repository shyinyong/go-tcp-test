package main

import (
	"github.com/shyinyong/go-tcp-test/config"
	"github.com/shyinyong/go-tcp-test/db/mysql"
	"github.com/shyinyong/go-tcp-test/db/redis"
	"github.com/shyinyong/go-tcp-test/handler/login"
	"log"
)

func main() {
	// Config env initialize
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
		return
	}
	// Initialize database connections
	store := mysql.NewDB(&cfg)
	// Create a Redis client
	redisClient := redis.NewClient()
	loginServer := login.NewServer(cfg, store, redisClient)
	loginServer.Start("localhost:8081") // Change the port as needed
}
