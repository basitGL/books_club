package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/basitGL/books_club/utils"
)

type RequestPayload struct {
	Rating   float32 `json:"rating"`
	BookID   int     `json:"book_id"`
	AuthorID int     `json:"author_id"`
}

type ResponsePayload struct {
	BookID        int     `json:"book_id"`
	UpdatedRating float32 `json:"rating"`
}

func RateBook(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
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

	responsePayload := ResponsePayload{
		BookID:        payload.BookID,
		UpdatedRating: updatedRating,
	}

	utils.SendResponse(w, "Rating updated successfully", "success", responsePayload, http.StatusOK)
}
