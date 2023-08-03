package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shyinyong/go-tcp-test/config"
	"github.com/shyinyong/go-tcp-test/db/mysql"
	"github.com/shyinyong/go-tcp-test/db/redis"
	"github.com/shyinyong/go-tcp-test/handler"
	"net/http"
	"os"
)

func main() {
	// Config env initialize
	cfg, err := config.LoadConfig(".")
	checkErr(err)
	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	mysqlDB, err := mysql.NewMySQLDB(&cfg)
	if err != nil {
		log.Print("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	redisDB, err := redis.NewRedisDB(&cfg)
	if err != nil {
		log.Print("Failed to connect to Redis: %v", err)
	}
	defer redisDB.Close()

	gameHandler := handler.NewGameHandler(mysqlDB, redisDB)
	// Set up HTTP routes and WebSocket handling using gameHandler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// WebSocket handling logic using gameHandler
		// 升级连接为 WebSocket
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		// 在这里你可以使用 gameHandler 处理 WebSocket 连接和游戏逻辑
		gameHandler.HandleWebSocket(conn)
	})

	// Start the HTTP server
	log.Print("Starting game server...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print("Failed to start server: %v", err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal()
	}
}
