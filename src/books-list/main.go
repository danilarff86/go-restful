package main

import (
	"log"
	"net/http"

	"books-list/controllers"
	"books-list/driver"

	"github.com/gorilla/mux"
)

func main() {
	controller := controllers.CreateController(driver.ConnectDB())

	router := mux.NewRouter()

	router.HandleFunc("/books", controller.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", controller.GetBook).Methods("GET")
	router.HandleFunc("/books", controller.AddBook).Methods("POST")
	router.HandleFunc("/books", controller.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", controller.RemoveBook).Methods("DELETE")

	log.Fatalln(http.ListenAndServe(":8000", router))
}
