package psql

import (
	"log"

	"books-list/models"

	"database/sql"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Queries struct {
	db *sql.DB
}

func CreateQueries(db *sql.DB) *Queries {
	return &Queries{db: db}
}

func (this *Queries) GetBooks() []models.Book {
	rows, err := this.db.Query("SELECT * FROM books")
	logFatal(err)

	defer rows.Close()

	var book models.Book
	books := []models.Book{}

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}

	return books
}

func (this *Queries) GetBook(id string) models.Book {
	rows, err := this.db.Query("SELECT * FROM books WHERE id=$1", id)
	logFatal(err)

	defer rows.Close()

	var book models.Book
	if rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
	}

	return book
}

func (this *Queries) AddBook(book *models.Book) int {
	var id int

	err := this.db.QueryRow("INSERT INTO books (title, author, year) VALUES($1, $2, $3) RETURNING id;",
		book.Title, book.Author, book.Year).Scan(&id)
	logFatal(err)

	return id
}

func (this *Queries) UpdateBook(book *models.Book) int64 {
	result, err := this.db.Exec("UPDATE books SET title=$1, author=$2, year=$3 WHERE id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &book.ID)
	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return rowsUpdated
}

func (this *Queries) RemoveBook(id string) int64 {
	result, err := this.db.Exec("DELETE FROM books WHERE id=$1 RETURNING id", &id)
	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return rowsUpdated
}
