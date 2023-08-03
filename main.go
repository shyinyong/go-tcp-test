package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-tcp-test/db"
	"go-tcp-test/db/enitty"
	"go-tcp-test/db/store"
	"go-tcp-test/utils"
	"os"
)

// App represents the application handlers
type App struct {
	db *db.DB
	//	UserService service.UserService
}

// NewApp creates a new App instance with injected DB
func NewApp(db *db.DB) *App {
	return &App{db: db}
}

// ListBooks handles listing all books
func (a *App) ListBooks() ([]enitty.Book, error) {
	return a.db.GetBooks()
}

func main() {
	// Config env initialize
	config, err := utils.LoadConfig(".")
	checkErr(err)
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Redis
	rc, err := store.NewRedisClient(
		config.RedisAddress,
		"",
		0,
	)
	if err != nil {
		log.Print(err)
	}

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379", // Redis server address
	//	Password: "",               // Password (if any)
	//	DB:       0,                // Database index
	//})

	// DB
	db, err := db.NewDB(config.DBSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// New App
	app := NewApp(db)
	// Example: List all books
	books, err := app.ListBooks()
	if err != nil {
		log.Print(err)
	}

	for _, book := range books {
		fmt.Printf("Book Isbn: %s, Author: %s\n", book.Isbn, book.Author)
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatal()
	}
}
