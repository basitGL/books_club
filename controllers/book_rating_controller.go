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

	defer tx.Rollback() // Ensure rollback if something goes wrong

	// Record exists, so update the rating
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

// func RateBook(w http.ResponseWriter, r *http.Request) {
// 	var payload RequestPayload
// 	var currentRating float32
// 	var updatedRating float32
// 	var ratingCount int
// 	var authorRating float32
// 	var updatedAuthorRating float32
// 	var authorTotalRating float32

// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Println(payload.BookID)

// 	tx, err := database.Begin()
// 	if err != nil {
// 		fmt.Println("starting transaction: error")
// 		http.Error(w, "Unable to start transaction", http.StatusInternalServerError)
// 		return
// 	}

// 	defer tx.Rollback() // Ensure rollback if something goes wrong

// 	// Attempt to get the current rating
// 	err = tx.QueryRow(`SELECT rating, rating_count FROM book_ratings WHERE book_id = ?`, payload.BookID).Scan(&currentRating, &ratingCount)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ratingCount = 1
// 			updatedRating = payload.Rating
// 			// No record found, so insert a new one with the initial rating
// 			insertFirstBookRating(w, tx, &payload, updatedRating, ratingCount)
// 			// _, err = tx.Exec(`INSERT INTO book_ratings (book_id, rating, rating_count) VALUES (?, ?)`, payload.BookID, updatedRating, ratingCount)
// 			// if err != nil {
// 			// 	fmt.Println("unable to insert rating:", err)
// 			// 	http.Error(w, "Unable to insert rating", http.StatusInternalServerError)
// 			// 	return
// 			// }
// 		} else {
// 			// An error other than ErrNoRows occurred
// 			fmt.Println("error querying current rating:", err)
// 			http.Error(w, "Database error", http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		// Record exists, so update the rating
// 		ratingCount++
// 		updatedRating = (payload.Rating + currentRating) / float32(ratingCount)
// 		_, err := tx.Exec(`UPDATE book_ratings SET rating = ?, rating_count = ? WHERE book_id = ?`, updatedRating, ratingCount, payload.BookID)
// 		if err != nil {
// 			fmt.Println("unable to update rating:", err)
// 			http.Error(w, "Unable to update rating", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	err = tx.QueryRow(`SELECT average_rating, total_ratings FROM author_ratings WHERE author_id = ?`, payload.AuthorID).Scan(&authorRating, &authorTotalRating)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			insertFirstAuthorRating(w, tx, &payload, updatedRating, ratingCount)
// 		}
// 	} else {
// 		authorTotalRating++
// 		updatedAuthorRating = (payload.Rating + authorRating) / float32(authorTotalRating)

// 		_, err := tx.Exec(`UPDATE author_ratings SET average_rating = ?, total_ratings = ? WHERE author_id = ?`, updatedAuthorRating, authorTotalRating, payload.AuthorID)
// 		if err != nil {
// 			fmt.Println("unable to update author rating:", err)
// 			http.Error(w, "Unable to update author rating", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	if err = tx.Commit(); err != nil {
// 		fmt.Println("committing transaction: error")
// 		http.Error(w, "Transaction commit error", http.StatusInternalServerError)
// 		return
// 	}

// 	responsePayload := ResponsePayload{
// 		BookID:        payload.BookID,
// 		UpdatedRating: updatedRating,
// 	}

// 	utils.SendResponse(w, "Rating updated successfully", "success", responsePayload, http.StatusOK)
// }

// func insertFirstBookRating(w http.ResponseWriter, tx *sql.Tx, payload *RequestPayload, updatedRating float32, ratingCount int) {
// 	_, err := tx.Exec(`INSERT INTO book_ratings (book_id, rating, rating_count) VALUES (?, ?, ?)`, payload.BookID, updatedRating, ratingCount)
// 	if err != nil {
// 		fmt.Println("unable to insert rating:", err)
// 		http.Error(w, "Unable to insert rating", http.StatusInternalServerError)
// 		return
// 	}
// }

// func insertFirstAuthorRating(w http.ResponseWriter, tx *sql.Tx, payload *RequestPayload, updatedRating float32, ratingCount int) {
// 	_, err := tx.Exec(`INSERT INTO author_ratings (author_id, average_rating, total_ratings) VALUES (?, ?, ?)`, payload.AuthorID, updatedRating, ratingCount)
// 	if err != nil {
// 		fmt.Println("unable to insert author rating:", err)
// 		http.Error(w, "Unable to insert author rating", http.StatusInternalServerError)
// 		return
// 	}
// }
