package routes

import (
	"net/http"

	"github.com/basitGL/books_club/controllers"
	"github.com/basitGL/books_club/services"
	"github.com/basitGL/books_club/utils"
	"github.com/gorilla/mux"
)

const uploadDir = "uploads"

type Router struct {
	router      *mux.Router
	authService *services.AuthService
}

func NewRouter(authService *services.AuthService) *Router {
	return &Router{
		router:      mux.NewRouter(),
		authService: authService,
	}
}

func (r *Router) Init() *mux.Router {

	print(http.Dir(uploadDir))

	r.router.Handle("/uploads/", http.StripPrefix("/uploads", http.FileServer(http.Dir(uploadDir))))

	// Create controllers with auth service
	userController := controllers.NewUserController(r.authService)
	bookController := controllers.NewBookController(r.authService)
	authorController := controllers.NewAuthorController(r.authService)

	// Public routes
	r.router.HandleFunc("/auth/register", userController.CreateUser).Methods("POST")
	r.router.HandleFunc("/auth/login", userController.LoginUser).Methods("POST")

	// Protected routes
	protected := r.router.PathPrefix("").Subrouter()
	protected.Use(utils.AuthMiddleware(r.authService))

	protected.HandleFunc("/books", bookController.GetAllBooks).Methods("GET")
	protected.HandleFunc("/book", bookController.AddBook).Methods("POST")
	protected.HandleFunc("/author", authorController.CreateAuthor).Methods("POST")
	protected.HandleFunc("/author/{id}", authorController.GetAuthor).Methods("GET")
	protected.HandleFunc("/rate-book", bookController.RateBook).Methods("POST")

	// Admin routes
	adminRouter := r.router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(utils.AuthMiddleware(r.authService))
	adminRouter.Use(utils.RoleMiddleware("admin"))
	// Add admin routes here

	return r.router
}
