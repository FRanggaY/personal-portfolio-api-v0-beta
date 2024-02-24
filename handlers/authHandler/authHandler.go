package authHandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {

}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserForm true "User input"
// @Success 201 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	// define input from json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		log.Fatal("Error decoding new user: ")
	}
	defer r.Body.Close()

	// hash
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)

	userInput.Password = string(hashPassword)

	// insert to database
	if err := models.DB.Create(&userInput).Error; err != nil {
		log.Fatal("Error creating new user")
	}

	response, _ := json.Marshal(map[string]string{"message": "success"})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {

}
