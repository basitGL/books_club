package utils

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, message string, status string, data interface{}, httpStatus int) {
	response := GeneralResponse{
		Message: message,
		Status:  status,
		Data:    data,
	}

	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}
