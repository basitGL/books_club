package models

import (
	"time"
)

type Book struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
	CoverPicture    string    `json:"cover_picture"`
	Price           float32   `json:"price"`
	AverageRating   float32   `json:"average_rating"`
	TotalRatings    int       `json:"total_ratings"`
	PublicationDate time.Time `json:"publication_date"`
	Author          Author    `json:"author"`
}
