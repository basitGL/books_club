package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/basitGL/books_club/models"
	"github.com/basitGL/books_club/utils"
)

var (
	authorId int
	name string
	description string
	image string
)

func CreateAuthor(w http.ResponseWriter, r *http.Request) {
	name = r.FormValue("name")
    description = r.FormValue("description")

	file, _, err := r.FormFile("avatar")
	if err!=nil {
		fmt.Println("Invalid File:")
		utils.SendResponse(w,"Invalid image file", "error", nil , http.StatusBadRequest)

		return
	}
	
	image, err = utils.UploadFileToServer(file)
	if err!=nil {
		fmt.Println(err.Error())
		utils.SendResponse(w,"Unable to upload image", "error", nil , http.StatusBadRequest)
		return;
	}


    result, err := database.Exec(`
        INSERT INTO authors (name, avatar, description) VALUES (?, ?, ?)
    `, name, image, description)

    if err != nil {
		fmt.Println("Failed to add book:");
		utils.SendResponse(w,"Unable to insert data", "error", nil , http.StatusBadRequest)
		return
	}

	authorId,_ := result.LastInsertId()

	res := map[string]interface{} {
    	"author_id": authorId,
	}

	utils.SendResponse(w,"Author created successfully.","success",res,http.StatusCreated)

}

func GetAuthor(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/");
	id ,err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		fmt.Println("Invalid author id")
		utils.SendResponse(w,"Invalid author id", "error", nil , http.StatusBadRequest)
		return
	}

	result := database.QueryRow(`SELECT id, name, avatar, description from authors WHERE id = ?`, id).Scan(&authorId, &name, &image, &description)
	if result != nil {
		if result == sql.ErrNoRows{
			fmt.Println("No records found")
			utils.SendResponse(w,"No author found for given id", "error", nil , http.StatusNotFound)
			return
		}
	}

	author := models.Author{
		ID: id,
		Name: name,
		Avatar: image,
		Desc: description,
	}

	utils.SendResponse(w,"Author found successfully", "success", author , http.StatusOK)
}