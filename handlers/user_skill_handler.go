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

// CreateUserSkill godoc
// @Summary Create a new user skill
// @Description Create a new user skill
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserSkillCreateForm true "User skill input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-skill [post]
func CreateUserSkill(w http.ResponseWriter, r *http.Request) {

	// define input from json
	var userSkillInput models.UserSkillCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userSkillInput); err != nil {
		log.Fatal("Error decoding new user skill: ")
	}
	defer r.Body.Close()

	userRepo := repositories.NewUserRepository()
	skillRepo := repositories.NewSkillRepository()
	userSkillRepo := repositories.NewUserSkillRepository()

	// validate user id
	_, user_err := userRepo.Read(userSkillInput.UserID)
	if user_err != nil {
		// Handle error
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate skill id
	_, skill_err := skillRepo.Read(userSkillInput.SkillID)
	if skill_err != nil {
		// Handle error
		response := map[string]string{"message": "Skill ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user id and skill id
	exist_data, _ := userSkillRepo.ReadByUserIDSkillID(userSkillInput.UserID, userSkillInput.SkillID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Skill already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newUserSkillData := models.UserSkill{
		UserID:  uint(userSkillInput.UserID),
		SkillID: uint(userSkillInput.SkillID),
	}

	// insert to database
	if newUserSkill, err := userSkillRepo.Create(&newUserSkillData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user skill"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserSkill.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user skill godoc
// @Summary Delete User Skill
// @Description Delete user skill
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Skill ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-skill/{id} [delete]
func DeleteUserSkill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userSkillIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Skill not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userSkillRepo := repositories.NewUserSkillRepository()
	userSkillID := helper.ParseIDStringToInt(userSkillIDStr)
	err := userSkillRepo.Delete(userSkillID)
	if err != nil {
		response := map[string]string{"message": "User Skill ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
