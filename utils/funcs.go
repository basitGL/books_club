package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads"

func UploadFileToServer(file multipart.File) (string, error) {

	defer file.Close()

	// Ensure the upload directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Create a new file in the upload directory
	filePath := filepath.Join(uploadDir, "uploaded_file") // Use the original filename or generate a new one
	destFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file on server: %w", err)
	}
	defer destFile.Close()

	// Copy the uploaded file data to the new file
	if _, err := io.Copy(destFile, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil // Return the file path and nil for no error
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
