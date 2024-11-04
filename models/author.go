package models

type Author struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Avatar        string  `json:"avatar"`
	Desc          string  `json:"description"`
	AverageRating float32 `json:"average_rating"`
	TotalRatings  int     `json:"total_ratings"`
}
