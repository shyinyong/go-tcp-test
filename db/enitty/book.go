package enitty

type Book struct {
	Isbn   string
	Title  string
	Author string
	Price  string
}

// GetBooks fetches all books from the database
//func (db *DB) GetBooks() ([]Book, error) {
//	rows, err := db.conn.Query("SELECT isbn, title, author, price FROM books")
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var books []Book
//	for rows.Next() {
//		book := Book{}
//		err := rows.Scan(&book.Isbn, &book.Title, &book.Author, &book.Price)
//		if err != nil {
//			panic(err)
//		}
//		books = append(books, book)
//	}
//	return books, nil
//}
