package main

import (
	"encoding/json"
	"log"
	"net/http"

	"database/sql"
	"github.com/lib/pq"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book
var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	pgURL, err := pq.ParseURL("postgres://postgres:docker@localhost/postgres?sslmode=disable")
	logFatal(err)

	db, err = sql.Open("postgres", pgURL)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatalln(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all books")
	var book Book
	books = []Book{}
	rows, err := db.Query("SELECT * FROM books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	log.Println("Get one book with id", id)

	rows, err := db.Query("SELECT * FROM books WHERE id=$1", id)
	logFatal(err)

	defer rows.Close()

	var book Book
	if rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
	}

	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Add one book")
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("INSERT INTO books (title, author, year) VALUES($1, $2, $3) RETURNING id;",
		book.Title, book.Author, book.Year).Scan(&book.ID)
	logFatal(err)

	json.NewEncoder(w).Encode(book.ID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Update a book")
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("UPDATE books SET title=$1, author=$2, year=$3 WHERE id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &book.ID)
	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	log.Println("Remove a book with id", id)

	result, err := db.Exec("DELETE FROM books WHERE id=$1 RETURNING id", &id)
	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}
