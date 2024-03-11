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

// Create User Project Translation godoc
// @Summary Create a new User Project Translation
// @Description Create a new User Project Translation with file upload support
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserProjectTranslationCreateForm true "User Project input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-project-translation [post]
func CreateUserProjectTranslation(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	// define input from json
	var userProjectTranslationInput models.UserProjectTranslationCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userProjectTranslationInput); err != nil {
		log.Fatal("Error decoding new user project: ")
	}
	defer r.Body.Close()

	userProjectRepo := repositories.NewUserProjectRepository()
	userProjectTranslationRepo := repositories.NewUserProjectTranslationRepository()

	// validate user_project id
	userProject, user_project_err := userProjectRepo.Read(userProjectTranslationInput.UserProjectID)
	if user_project_err != nil {
		// Handle error
		response := map[string]string{"message": "Project not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if userProject.UserID != uint(userID) {
		response := map[string]string{"message": "Not allowing create"}
		helper.ResponseJSON(w, http.StatusForbidden, response)
		return
	}

	exist_data, _ := userProjectTranslationRepo.ReadByLanguageIDUserProjectID(userProjectTranslationInput.LanguageID, userProjectTranslationInput.UserProjectID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Project already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// Create a new UserProjectTranslation record
	newUserProjectTranslationData := models.UserProjectTranslation{
		LanguageID:    uint(userProjectTranslationInput.LanguageID),
		UserProjectID: uint(userProjectTranslationInput.UserProjectID),
		Name:          userProjectTranslationInput.Name,
		Description:   userProjectTranslationInput.Description,
	}
	// insert to database
	if newUserProjectTranslation, err := userProjectTranslationRepo.Create(&newUserProjectTranslationData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user project translation"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserProjectTranslation.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user project translation godoc
// @Summary Delete User project translation
// @Description Delete user project translation
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Project Translation ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Success 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-project-translation/{id} [delete]
func DeleteUserProjectTranslation(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	userProjectTranslationIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Project Translation not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userProjectRepo := repositories.NewUserProjectRepository()
	userProjectTranslationRepo := repositories.NewUserProjectTranslationRepository()
	userProjectTranslationID := helper.ParseIDStringToInt(userProjectTranslationIDStr)

	userProjectTranslation, err := userProjectTranslationRepo.Read(userProjectTranslationID)
	if err != nil {
		response := map[string]string{"message": "User Project Translation ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	userProject, err := userProjectRepo.Read(int64(userProjectTranslation.UserProjectID))
	if err != nil {
		response := map[string]string{"message": "User Project ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	if userProject.UserID != uint(userID) {
		response := map[string]string{"message": "Not allowing delete"}
		helper.ResponseJSON(w, http.StatusForbidden, response)
		return
	}

	errDelete := userProjectTranslationRepo.Delete(userProjectTranslationID)
	if errDelete != nil {
		response := map[string]string{"message": "User Project Translation ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
