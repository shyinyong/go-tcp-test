package db

import (
	"database/sql"
	"github.com/shyinyong/go-tcp-test/db/enitty"
)
import _ "github.com/go-sql-driver/mysql"

// DB represents the database operations
type DB struct {
	conn *sql.DB
}

// NewDB creates a new DB instance
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{conn: db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// GetBooks fetches all books from the database
func (db *DB) GetBooks() ([]enitty.Book, error) {
	rows, err := db.conn.Query("SELECT isbn, title, author, price FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []enitty.Book
	for rows.Next() {
		book := enitty.Book{}
		err := rows.Scan(&book.Isbn, &book.Title, &book.Author, &book.Price)
		if err != nil {
			panic(err)
		}
		books = append(books, book)
	}
	return books, nil
}
