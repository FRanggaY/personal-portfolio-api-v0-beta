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

// Create User Education Translation godoc
// @Summary Create a new User Education Translation
// @Description Create a new User Education Translation with file upload support
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserEducationTranslationCreateForm true "User skill input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-education-translation [post]
func CreateUserEducationTranslation(w http.ResponseWriter, r *http.Request) {
	// define input from json
	var userEducationTranslationInput models.UserEducationTranslationCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userEducationTranslationInput); err != nil {
		log.Fatal("Error decoding new user skill: ")
	}
	defer r.Body.Close()

	languageRepo := repositories.NewLanguageRepository()
	userEducationRepo := repositories.NewUserEducationRepository()
	userEducationTranslationRepo := repositories.NewUserEducationTranslationRepository()

	// validate language id
	_, language_err := languageRepo.Read(userEducationTranslationInput.LanguageID)
	if language_err != nil {
		// Handle error
		response := map[string]string{"message": "Language ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user_education id
	_, user_education_err := userEducationRepo.Read(userEducationTranslationInput.UserEducationID)
	if user_education_err != nil {
		// Handle error
		response := map[string]string{"message": "Skill ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	exist_data, _ := userEducationTranslationRepo.ReadByLanguageIdUserEducationId(userEducationTranslationInput.LanguageID, userEducationTranslationInput.UserEducationID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Education already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Create a new UserEducationTranslation record
	newUserEducationTranslationData := models.UserEducationTranslation{
		LanguageID:      uint(userEducationTranslationInput.LanguageID),
		UserEducationID: uint(userEducationTranslationInput.UserEducationID),
		Title:           userEducationTranslationInput.Title,
		Description:     userEducationTranslationInput.Description,
		Category:        userEducationTranslationInput.Category,
		Location:        userEducationTranslationInput.Location,
		LocationType:    userEducationTranslationInput.LocationType,
		MonthStart:      userEducationTranslationInput.MonthStart,
		MonthEnd:        userEducationTranslationInput.MonthEnd,
		YearStart:       userEducationTranslationInput.YearStart,
		YearEnd:         userEducationTranslationInput.YearEnd,
	}
	// insert to database
	if newUserEducationTranslation, err := userEducationTranslationRepo.Create(&newUserEducationTranslationData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user education translation"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserEducationTranslation.Id,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user Education Translation godoc
// @Summary Delete User Education Translation
// @Description Delete user Education Translation
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Education Translation ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-education-translation/{id} [delete]
func DeleteUserEducationTranslation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userEducationTranslationIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Education Translation not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userEducationTranslationRepo := repositories.NewUserEducationTranslationRepository()
	userEducationTranslationID := helper.ParseIDStringToInt(userEducationTranslationIDStr)
	err := userEducationTranslationRepo.Delete(userEducationTranslationID)
	if err != nil {
		response := map[string]string{"message": "User Education Translation ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
