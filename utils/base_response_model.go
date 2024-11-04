package utils

type GeneralResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"` // Use interface{} to allow flexibility for different data types
}