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

// CreateProjectPlatformTranslation godoc
// @Summary Create a new project Platform translation
// @Description Create a new project Platform translation
// @Tags project-platforms
// @Accept json
// @Produce json
// @Param input body models.ProjectPlatformTranslationCreateForm true "Project Platform translation input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /project-platform-translation [post]
func CreateProjectPlatformTranslation(w http.ResponseWriter, r *http.Request) {

	// define input from json
	var projectPlatformTranslationInput models.ProjectPlatformTranslationCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&projectPlatformTranslationInput); err != nil {
		log.Fatal("Error decoding new project platform translation: ")
	}
	defer r.Body.Close()

	languageRepo := repositories.NewLanguageRepository()
	projectPlatformRepo := repositories.NewProjectPlatformRepository()
	projectPlatformTranslationRepo := repositories.NewProjectPlatformTranslationRepository()

	// validate lang id
	_, lang_err := languageRepo.Read(projectPlatformTranslationInput.LanguageID)
	if lang_err != nil {
		// Handle error
		response := map[string]string{"message": "Lang ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate projectPlatform id
	_, projectPlatform_err := projectPlatformRepo.Read(projectPlatformTranslationInput.ProjectPlatformID)
	if projectPlatform_err != nil {
		// Handle error
		response := map[string]string{"message": "Project Platform ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate lang id and projectPlatform id
	exist_data, _ := projectPlatformTranslationRepo.ReadByLanguageIDProjectPlatformID(projectPlatformTranslationInput.LanguageID, projectPlatformTranslationInput.ProjectPlatformID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Project Platform Translation already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newProjectPlatformTranslationData := models.ProjectPlatformTranslation{
		LanguageID:        uint(projectPlatformTranslationInput.LanguageID),
		ProjectPlatformID: uint(projectPlatformTranslationInput.ProjectPlatformID),
		Title:             projectPlatformTranslationInput.Title,
		Description:       projectPlatformTranslationInput.Description,
	}

	// insert to database
	if newProjectPlatformTranslation, err := projectPlatformTranslationRepo.Create(&newProjectPlatformTranslationData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new project platform translation"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newProjectPlatformTranslation.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user projectPlatform godoc
// @Summary Delete User Project Platform
// @Description Delete user project Platform
// @Tags project-platforms
// @Accept json
// @Produce json
// @Param id path int true "Project Platform Translation ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /project-platform-translation/{id} [delete]
func DeleteProjectPlatformTranslation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectPlatformTranslationIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "Project Platform Translation not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	projectPlatformTranslationRepo := repositories.NewProjectPlatformTranslationRepository()
	projectPlatformTranslationID := helper.ParseIDStringToInt(projectPlatformTranslationIDStr)
	err := projectPlatformTranslationRepo.Delete(projectPlatformTranslationID)
	if err != nil {
		response := map[string]string{"message": "Project Platform Translation ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
