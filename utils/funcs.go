package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/basitGL/books_club/services"
	"github.com/golang-jwt/jwt"
)

const uploadDir = "./uploads"

// uploads a picture to server
func UploadFileToServer(file multipart.File, handler *multipart.FileHeader, r *http.Request) (string, error) {
	defer file.Close()

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	filename := generateFilename(handler.Filename)

	destFile, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return "", fmt.Errorf("failed to create file on server: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	imageURL := fmt.Sprintf("http://%s/uploads/%s", r.Host, filename)
	return imageURL, nil
}

func generateFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
}

// middleware for setting content-type
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// middleware for authorization
func AuthMiddleware(authService *services.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			token, err := authService.ValidateToken(tokenString)
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// role-based authorization middleware
func RoleMiddleware(requiredRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("user").(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userRole, ok := claims["role"].(string)
			if !ok || userRole != requiredRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
