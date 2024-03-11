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

// Create User Language Translation godoc
// @Summary Create a new User Language Translation
// @Description Create a new User Language Translation with file upload support
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserLanguageTranslationCreateForm true "User Language input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-language-translation [post]
func CreateUserLanguageTranslation(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	// define input from json
	var userLanguageTranslationInput models.UserLanguageTranslationCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userLanguageTranslationInput); err != nil {
		log.Fatal("Error decoding new user language: ")
	}
	defer r.Body.Close()

	languageRepo := repositories.NewLanguageRepository()
	userLanguageRepo := repositories.NewUserLanguageRepository()
	userLanguageTranslationRepo := repositories.NewUserLanguageTranslationRepository()

	// validate language id
	_, language_err := languageRepo.Read(userLanguageTranslationInput.LanguageID)
	if language_err != nil {
		// Handle error
		response := map[string]string{"message": "Language ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user_language id
	userLanguage, user_language_err := userLanguageRepo.ReadByUserIDLanguageID(userID, userLanguageTranslationInput.SelectLanguageID)
	if user_language_err != nil {
		// Handle error
		response := map[string]string{"message": "Select Language not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	exist_data, _ := userLanguageTranslationRepo.ReadByLanguageIDUserLanguageID(userLanguageTranslationInput.LanguageID, userLanguage.ID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Language already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Create a new UserLanguageTranslation record
	newUserLanguageTranslationData := models.UserLanguageTranslation{
		LanguageID:     uint(userLanguageTranslationInput.LanguageID),
		UserLanguageID: uint(userLanguage.ID),
		Description:    userLanguageTranslationInput.Description,
		Title:          userLanguageTranslationInput.Title,
	}
	// insert to database
	if newUserLanguageTranslation, err := userLanguageTranslationRepo.Create(&newUserLanguageTranslationData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user language translation"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserLanguageTranslation.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user Language Translation godoc
// @Summary Delete User Language Translation
// @Description Delete user Language Translation
// @Tags users
// @Accept json
// @Produce json
// @Param select_language_id path int true "Select Language ID"
// @Param language_id path int true "language ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-language-translation/{select_language_id}/{language_id} [delete]
func DeleteUserLanguageTranslation(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	languageIDStr, ok := vars["language_id"]
	if !ok {
		response := map[string]string{"message": "Language ID not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	selectLanguageIDStr, ok := vars["select_language_id"]
	if !ok {
		response := map[string]string{"message": "Select Language ID not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	languageRepo := repositories.NewLanguageRepository()
	userLanguageRepo := repositories.NewUserLanguageRepository()
	userLanguageTranslationRepo := repositories.NewUserLanguageTranslationRepository()
	languageID := helper.ParseIDStringToInt(languageIDStr)
	selectLanguageID := helper.ParseIDStringToInt(selectLanguageIDStr)

	_, language_err := languageRepo.Read(selectLanguageID)
	if language_err != nil {
		// Handle error
		response := map[string]string{"message": "Select Language not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	userLanguage, user_language_err := userLanguageRepo.ReadByUserIDLanguageID(userID, selectLanguageID)
	if user_language_err != nil {
		// Handle error
		response := map[string]string{"message": "Language not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	err := userLanguageTranslationRepo.DeleteByLanguageIDUserLanguageID(languageID, userLanguage.ID)
	if err != nil {
		response := map[string]string{"message": "User Language Translation ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
