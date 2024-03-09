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

// CreateSkillTranslation godoc
// @Summary Create a new skill translation
// @Description Create a new skill translation
// @Tags skills
// @Accept json
// @Produce json
// @Param input body models.SkillTranslationCreateForm true "Skill translation input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /skill-translation [post]
func CreateSkillTranslation(w http.ResponseWriter, r *http.Request) {

	// define input from json
	var skillTranslationInput models.SkillTranslationCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&skillTranslationInput); err != nil {
		log.Fatal("Error decoding new skill translation: ")
	}
	defer r.Body.Close()

	languageRepo := repositories.NewLanguageRepository()
	skillRepo := repositories.NewSkillRepository()
	skillTranslationRepo := repositories.NewSkillTranslationRepository()

	// validate lang id
	_, lang_err := languageRepo.Read(skillTranslationInput.LanguageID)
	if lang_err != nil {
		// Handle error
		response := map[string]string{"message": "Lang ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate skill id
	_, skill_err := skillRepo.Read(skillTranslationInput.SkillID)
	if skill_err != nil {
		// Handle error
		response := map[string]string{"message": "Skill ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate lang id and skill id
	exist_data, _ := skillTranslationRepo.ReadByLanguageIdSkillId(skillTranslationInput.LanguageID, skillTranslationInput.SkillID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Skill Translation already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newSkillTranslationData := models.SkillTranslation{
		LanguageID:  uint(skillTranslationInput.LanguageID),
		SkillID:     uint(skillTranslationInput.SkillID),
		Description: skillTranslationInput.Description,
	}

	// insert to database
	if newSkillTranslation, err := skillTranslationRepo.Create(&newSkillTranslationData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new skill translation"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newSkillTranslation.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user skill godoc
// @Summary Delete User Skill
// @Description Delete user skill
// @Tags skills
// @Accept json
// @Produce json
// @Param id path int true "Skill Translation ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /skill-translation/{id} [delete]
func DeleteSkillTranslation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skillTranslationIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "Skill Translation not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	skillTranslationRepo := repositories.NewSkillTranslationRepository()
	skillTranslationID := helper.ParseIDStringToInt(skillTranslationIDStr)
	err := skillTranslationRepo.Delete(skillTranslationID)
	if err != nil {
		response := map[string]string{"message": "Skill Translation ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
