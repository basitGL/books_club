package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/basitGL/books_club/models"
	"github.com/basitGL/books_club/services"
	"github.com/basitGL/books_club/utils"
)

type UserController struct {
	authService *services.AuthService
}

func NewUserController(authService *services.AuthService) *UserController {
	return &UserController{
		authService: authService,
	}
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserImage string `json:"user_image"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash password using auth service
	hashedPassword, err := c.authService.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Create user object
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserImage: req.UserImage,
	}

	_, err = database.Exec(`INSERT INTO users (first_name, last_name, email, password, user_image, role) VALUES (?,?,?,?,?,?)`, user.FirstName, user.LastName, user.Email, user.Password, user.UserImage, user.Role)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	tokenPair, err := c.authService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}

// TODO: Implement login and other user controller methods

func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginUserRequest
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		utils.SendResponse(w, "Invalid request payload", "error", nil, http.StatusBadRequest)
		return
	}

	statement := database.QueryRow(`SELECT id, first_name, last_name, email, password, user_image, role  from users where email = ?`, loginReq.Email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.UserImage,
		&user.Role,
	)
	if statement != nil {
		if statement == sql.ErrNoRows {
			utils.SendResponse(w, "No user found for given email", "error", nil, http.StatusNotFound)
			return
		}

	}

	if !c.authService.CheckPassword(loginReq.Password, user.Password) {
		utils.SendResponse(w, "Incorrect password for provided email", "error", nil, http.StatusNotFound)
		return
	}

	token, err := c.authService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": token}

	utils.SendResponse(w, "Login Successful", "sucess", response, http.StatusOK)

}

// TODO: implement user controller methods

func EditUser() {

}

func DeleteUser() {

}
