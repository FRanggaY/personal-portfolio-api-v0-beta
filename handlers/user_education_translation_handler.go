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
// @Param input body models.UserEducationTranslationCreateForm true "User education input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-education-translation [post]
func CreateUserEducationTranslation(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	// define input from json
	var userEducationTranslationInput models.UserEducationTranslationCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userEducationTranslationInput); err != nil {
		print(err)
		log.Fatal("Error decoding new user education: ")
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
	userEducation, user_education_err := userEducationRepo.ReadByUserIDSchoolID(userID, userEducationTranslationInput.SchoolID)
	if user_education_err != nil {
		// Handle error
		response := map[string]string{"message": "Education not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	exist_data, _ := userEducationTranslationRepo.ReadByLanguageIDUserEducationID(userEducationTranslationInput.LanguageID, userEducation.ID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Education already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Create a new UserEducationTranslation record
	newUserEducationTranslationData := models.UserEducationTranslation{
		LanguageID:      uint(userEducationTranslationInput.LanguageID),
		UserEducationID: uint(userEducation.ID),
		Title:           userEducationTranslationInput.Title,
		Description:     userEducationTranslationInput.Description,
		Category:        userEducationTranslationInput.Category,
		Location:        userEducationTranslationInput.Location,
		LocationType:    userEducationTranslationInput.LocationType,
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
				"id": newUserEducationTranslation.ID,
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
// @Param school_id path int true "School ID"
// @Param language_id path int true "Language ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-education-translation/{school_id}/{language_id} [delete]
func DeleteUserEducationTranslation(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	schoolIDStr, ok := vars["school_id"]
	if !ok {
		response := map[string]string{"message": "School ID not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	languageIDStr, ok := vars["language_id"]
	if !ok {
		response := map[string]string{"message": "Language ID not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userEducationRepo := repositories.NewUserEducationRepository()
	userEducationTranslationRepo := repositories.NewUserEducationTranslationRepository()
	schoolID := helper.ParseIDStringToInt(schoolIDStr)
	languageID := helper.ParseIDStringToInt(languageIDStr)

	userEducation, user_education_err := userEducationRepo.ReadByUserIDSchoolID(userID, schoolID)
	if user_education_err != nil {
		// Handle error
		response := map[string]string{"message": "Education not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	err := userEducationTranslationRepo.DeleteByLanguageIDUserEducationID(languageID, userEducation.ID)
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
