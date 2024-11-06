package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/basitGL/books_club/config"
	"github.com/basitGL/books_club/models"
	"github.com/basitGL/books_club/services"
	"github.com/basitGL/books_club/utils"
)

var (
	database = config.Database()
)

type RateBookRequestPayload struct {
	Rating   float32 `json:"rating"`
	BookID   int     `json:"book_id"`
	AuthorID int     `json:"author_id"`
}

type RateBookResponsePayload struct {
	BookID        int     `json:"book_id"`
	UpdatedRating float32 `json:"rating"`
}

type BookController struct {
	authService *services.AuthService
}

func NewBookController(authService *services.AuthService) *BookController {
	return &BookController{
		authService: authService,
	}
}

func (c *BookController) GetAllBooks(w http.ResponseWriter, r *http.Request) {

	statement, err := database.Query(`
		SELECT 
		b.id AS book_id,
		b.title, 
		b.summary, 
		b.publication_date, 
		b.cover_picture, 
		b.price,
		IFNULL(br.rating, 0) AS book_average_rating, 
		IFNULL(br.rating_count, 0) AS book_total_ratings,
    
		a.id AS author_id,
		a.name AS author_name,
		a.avatar AS author_avatar,
		a.description AS author_description,
		
		IFNULL(ar.average_rating, 0) AS author_average_rating, 
		IFNULL(ar.total_ratings, 0) AS author_total_ratings

		FROM books b
		LEFT JOIN book_ratings br ON b.id = br.book_id
		LEFT JOIN book_authors ba ON b.id = ba.book_id
		LEFT JOIN authors a ON ba.author_id = a.id
		LEFT JOIN author_ratings ar ON a.id = ar.author_id;
	`)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
		return
	}

	defer statement.Close()

	var books []models.Book

	for statement.Next() {
		var book models.Book
		var publicationDateStr string

		err = statement.Scan(
			&book.ID,
			&book.Title,
			&book.Summary,
			&publicationDateStr,
			&book.CoverPicture,
			&book.Price,
			&book.AverageRating,
			&book.TotalRatings,
			&book.Author.ID,
			&book.Author.Name,
			&book.Author.Avatar,
			&book.Author.Desc,
			&book.Author.AverageRating,
			&book.Author.TotalRatings,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		book.PublicationDate, err = time.Parse("2006-01-02", publicationDateStr)
		if err != nil {
			log.Fatal("Failed to parse publication date:", err)
		}

		books = append(books, book)
	}

	utils.SendResponse(w, "Books retrieved successfully", "success", books, http.StatusOK)

}

func (c *BookController) AddBook(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	summary := r.FormValue("summary")
	publicationDate := r.FormValue("publication_date")
	authorId, err := strconv.Atoi(r.FormValue("author_id"))
	if err != nil {
		fmt.Println("Error parsing authorId value: ")
	}

	file, handler, err := r.FormFile("cover_picture")
	if err != nil {
		fmt.Println("Invalid File:")
	}

	coverPicture, err := utils.UploadFileToServer(file, handler, r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result, err := database.Exec(`
        INSERT INTO books (title, summary, publication_date, cover_picture) VALUES (?, ?, ?, ?)
    `, title, summary, publicationDate, coverPicture)

	if err != nil {
		fmt.Println("Failed to add book:")
	}

	bookId, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Failed to add book:")
	}

	_, err = database.Exec(`
        INSERT INTO book_authors (book_id, author_id) VALUES (?, ?)
    `, bookId, authorId)

	if err != nil {
		fmt.Println("Failed to add book:")
	}
	utils.SendResponse(w, "", "success", "", http.StatusCreated)
}

func (c *BookController) RateBook(w http.ResponseWriter, r *http.Request) {
	var payload RateBookRequestPayload
	var updatedRating float32

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := database.Begin()
	if err != nil {
		http.Error(w, "Unable to start transaction", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE book_ratings SET rating = ((rating * rating_count) + ?) / (rating_count + 1), rating_count = (rating_count + 1) WHERE book_id = ?`, payload.Rating, payload.BookID)
	if err != nil {
		http.Error(w, "Unable to update rating", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`UPDATE author_ratings SET average_rating = ((average_rating * total_ratings) + ?) / (total_ratings + 1), total_ratings = (total_ratings + 1) WHERE author_id = ?`, payload.Rating, payload.AuthorID)
	if err != nil {
		http.Error(w, "Unable to update author rating", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "Transaction commit error", http.StatusInternalServerError)
		return
	}

	responsePayload := RateBookResponsePayload{
		BookID:        payload.BookID,
		UpdatedRating: updatedRating,
	}

	utils.SendResponse(w, "Rating updated successfully", "success", responsePayload, http.StatusOK)
}
