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

// Create User Experience Translation godoc
// @Summary Create a new User Experience Translation
// @Description Create a new User Experience Translation with file upload support
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserExperienceTranslationCreateForm true "User experience input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-experience-translation [post]
func CreateUserExperienceTranslation(w http.ResponseWriter, r *http.Request) {
	// define input from json
	var userExperienceTranslationInput models.UserExperienceTranslationCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userExperienceTranslationInput); err != nil {
		log.Fatal("Error decoding new user experience: ")
	}
	defer r.Body.Close()

	languageRepo := repositories.NewLanguageRepository()
	userExperienceRepo := repositories.NewUserExperienceRepository()
	userExperienceTranslationRepo := repositories.NewUserExperienceTranslationRepository()

	// validate language id
	_, language_err := languageRepo.Read(userExperienceTranslationInput.LanguageID)
	if language_err != nil {
		// Handle error
		response := map[string]string{"message": "Language ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user_experience id
	_, user_experience_err := userExperienceRepo.Read(userExperienceTranslationInput.UserExperienceID)
	if user_experience_err != nil {
		// Handle error
		response := map[string]string{"message": "Skill ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	exist_data, _ := userExperienceTranslationRepo.ReadByLanguageIDUserExperienceID(userExperienceTranslationInput.LanguageID, userExperienceTranslationInput.UserExperienceID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Experience already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Create a new UserExperienceTranslation record
	newUserExperienceTranslationData := models.UserExperienceTranslation{
		LanguageID:       uint(userExperienceTranslationInput.LanguageID),
		UserExperienceID: uint(userExperienceTranslationInput.UserExperienceID),
		Title:            userExperienceTranslationInput.Title,
		Description:      userExperienceTranslationInput.Description,
		Category:         userExperienceTranslationInput.Category,
		Location:         userExperienceTranslationInput.Location,
		LocationType:     userExperienceTranslationInput.LocationType,
		Industry:         userExperienceTranslationInput.Industry,
		MonthStart:       userExperienceTranslationInput.MonthStart,
		MonthEnd:         userExperienceTranslationInput.MonthEnd,
		YearStart:        uint(userExperienceTranslationInput.YearStart),
		YearEnd:          uint(userExperienceTranslationInput.YearEnd),
	}
	// insert to database
	if newUserExperienceTranslation, err := userExperienceTranslationRepo.Create(&newUserExperienceTranslationData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user experience translation"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserExperienceTranslation.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user Experience Translation godoc
// @Summary Delete User Experience Translation
// @Description Delete user Experience Translation
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Experience Translation ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-experience-translation/{id} [delete]
func DeleteUserExperienceTranslation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userExperienceTranslationIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Experience Translation not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userExperienceTranslationRepo := repositories.NewUserExperienceTranslationRepository()
	userExperienceTranslationID := helper.ParseIDStringToInt(userExperienceTranslationIDStr)
	err := userExperienceTranslationRepo.Delete(userExperienceTranslationID)
	if err != nil {
		response := map[string]string{"message": "User Experience Translation ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
