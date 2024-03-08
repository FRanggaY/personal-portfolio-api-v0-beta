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

// CreateUserExperience godoc
// @Summary Create a new user experience
// @Description Create a new user experience
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserExperienceCreateForm true "User experience input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-experience [post]
func CreateUserExperience(w http.ResponseWriter, r *http.Request) {

	// define input from json
	var userExperienceInput models.UserExperienceCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userExperienceInput); err != nil {
		log.Fatal("Error decoding new user experience: ")
	}
	defer r.Body.Close()

	userRepo := repositories.NewUserRepository()
	companyRepo := repositories.NewCompanyRepository()
	userExperienceRepo := repositories.NewUserExperienceRepository()

	// validate user id
	_, user_err := userRepo.Read(userExperienceInput.UserId)
	if user_err != nil {
		// Handle error
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate company id
	_, company_err := companyRepo.Read(userExperienceInput.CompanyId)
	if company_err != nil {
		// Handle error
		response := map[string]string{"message": "Company ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user id and company id
	exist_data, _ := userExperienceRepo.ReadByUserIdCompanyId(userExperienceInput.UserId, userExperienceInput.CompanyId)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Experience already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newUserExperienceData := models.UserExperience{
		UserID:    uint(userExperienceInput.UserId),
		CompanyID: uint(userExperienceInput.CompanyId),
	}

	// insert to database
	if newUserExperience, err := userExperienceRepo.Create(&newUserExperienceData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user experience"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserExperience.Id,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user Experience godoc
// @Summary Delete User Experience
// @Description Delete user Experience
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Experience ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-experience/{id} [delete]
func DeleteUserExperience(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userExperienceIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Experience not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userExperienceRepo := repositories.NewUserExperienceRepository()
	userExperienceID := helper.ParseIDStringToInt(userExperienceIDStr)
	err := userExperienceRepo.Delete(userExperienceID)
	if err != nil {
		response := map[string]string{"message": "User Experience ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
