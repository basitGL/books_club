package routes

import (
	"github.com/basitGL/books_club/controllers"
	"github.com/gorilla/mux"
)

func Init() *mux.Router {
	route := mux.NewRouter()

	route.HandleFunc("/", controllers.GetAllBooks).Methods("GET")
	route.HandleFunc("/books", controllers.GetAllBooks).Methods("GET")
	route.HandleFunc("/book", controllers.AddBook).Methods("POST")
	route.HandleFunc("/author", controllers.CreateAuthor).Methods("POST")
	route.HandleFunc("/author/{id}", controllers.GetAuthor).Methods("GET")
	route.HandleFunc("/rate-book", controllers.RateBook).Methods("POST")

	return route
}
