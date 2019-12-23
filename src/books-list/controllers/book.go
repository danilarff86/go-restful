package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"books-list/models"
	"books-list/psql"

	"database/sql"

	"github.com/gorilla/mux"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Controller struct {
	queries *psql.Queries
}

func CreateController(db *sql.DB) *Controller {
	return &Controller{queries: psql.CreateQueries(db)}
}

func (this *Controller) GetBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all books")

	books := this.queries.GetBooks()

	json.NewEncoder(w).Encode(books)
}

func (this *Controller) GetBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	log.Println("Get one book with id", id)

	book := this.queries.GetBook(id)

	json.NewEncoder(w).Encode(book)
}

func (this *Controller) AddBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Add one book")
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)

	id := this.queries.AddBook(&book)

	json.NewEncoder(w).Encode(id)
}

func (this *Controller) UpdateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Update a book")
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)

	rowsUpdated := this.queries.UpdateBook(&book)

	json.NewEncoder(w).Encode(rowsUpdated)
}

func (this *Controller) RemoveBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	log.Println("Remove a book with id", id)

	rowsUpdated := this.queries.RemoveBook(id)

	json.NewEncoder(w).Encode(rowsUpdated)
}
