package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// CreateUserLanguage godoc
// @Summary Create a new user Language
// @Description Create a new user Language
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserLanguageCreateForm true "User Language input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-language [post]
func CreateUserLanguage(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	// define input from json
	var userLanguageInput models.UserLanguageCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userLanguageInput); err != nil {
		log.Fatal("Error decoding new user language: ")
	}
	defer r.Body.Close()

	userRepo := repositories.NewUserRepository()
	languageRepo := repositories.NewLanguageRepository()
	userLanguageRepo := repositories.NewUserLanguageRepository()

	// validate user id
	_, user_err := userRepo.Read(userID)
	if user_err != nil {
		// Handle error
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate language id
	_, lang_err := languageRepo.Read(userLanguageInput.LanguageID)
	if lang_err != nil {
		// Handle error
		response := map[string]string{"message": "Lang ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user id and language id
	exist_data, _ := userLanguageRepo.ReadByUserIDLanguageID(userID, userLanguageInput.LanguageID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Language already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newUserLanguageData := models.UserLanguage{
		UserID:     uint(userID),
		LanguageID: uint(userLanguageInput.LanguageID),
	}

	// insert to database
	if newUserLanguage, err := userLanguageRepo.Create(&newUserLanguageData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user language"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserLanguage.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user language godoc
// @Summary Delete User language
// @Description Delete user language
// @Tags users
// @Accept json
// @Produce json
// @Param language_id path int true "Language ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-language/{language_id} [delete]
func DeleteUserLanguage(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	languageIDStr, ok := vars["language_id"]
	if !ok {
		response := map[string]string{"message": "Language ID not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userLanguageRepo := repositories.NewUserLanguageRepository()
	languageID := helper.ParseIDStringToInt(languageIDStr)
	err := userLanguageRepo.DeleteByUserIDLanguageID(userID, languageID)
	if err != nil {
		response := map[string]string{"message": "User Language ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
