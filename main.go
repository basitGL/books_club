package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/basitGL/books_club/routes"
	"github.com/basitGL/books_club/services"
	"github.com/ichtrojan/thoth"
	"github.com/joho/godotenv"
)

func main() {

	logger, _ := thoth.Init("log")

	if err := godotenv.Load(); err != nil {
		logger.Log(errors.New("no .env file found"))
		log.Fatal("No .env file found")
	}

	port, exist := os.LookupEnv("PORT")

	if !exist {
		logger.Log(errors.New("PORT not set in .env"))
		log.Fatal("PORT not set in .env")
	}

	authService := services.NewAuthService("your-secret-key")
	router := routes.NewRouter(authService)
	r := router.Init()

	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started at http://localhost:" + port)
	err := http.ListenAndServe(":"+port, r)
	// err := http.ListenAndServe(":"+port, utils.ContentTypeMiddleware(r))

	if err != nil {
		logger.Log(err)
		log.Fatal(err)
	}
}
