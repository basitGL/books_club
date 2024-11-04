package routes

import (
	"github.com/basitGL/books_club/controllers"
	"github.com/gorilla/mux"
	// "github.com/basitGL/books_club/controllers"
)

func Init() *mux.Router {
	route := mux.NewRouter()

	route.HandleFunc("/", controllers.GetAllBooks).Methods("GET")
	route.HandleFunc("/books", controllers.GetAllBooks).Methods("GET")
	route.HandleFunc("/book", controllers.AddBook).Methods("POST")
	route.HandleFunc("/author", controllers.CreateAuthor).Methods("POST")
	route.HandleFunc("/author/{id}", controllers.GetAuthor).Methods("GET")
	route.HandleFunc("/rate-book", controllers.RateBook).Methods("POST")
	// route.HandleFunc("/delete/{id}", controllers.Delete)
	// route.HandleFunc("/complete/{id}", controllers.Complete)

	return route
}
