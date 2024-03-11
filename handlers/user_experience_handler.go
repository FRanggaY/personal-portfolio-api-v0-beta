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
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

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
	_, user_err := userRepo.Read(userID)
	if user_err != nil {
		// Handle error
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate company id
	_, company_err := companyRepo.Read(userExperienceInput.CompanyID)
	if company_err != nil {
		// Handle error
		response := map[string]string{"message": "Company ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user id and company id
	exist_data, _ := userExperienceRepo.ReadByUserIDCompanyID(userID, userExperienceInput.CompanyID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Experience already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newUserExperienceData := models.UserExperience{
		UserID:     uint(userID),
		CompanyID:  uint(userExperienceInput.CompanyID),
		MonthStart: userExperienceInput.MonthStart,
		MonthEnd:   userExperienceInput.MonthEnd,
		YearStart:  uint(userExperienceInput.YearStart),
		YearEnd:    uint(userExperienceInput.YearEnd),
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
				"id": newUserExperience.ID,
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
// @Param company_id path int true "Company ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-experience/{company_id} [delete]
func DeleteUserExperience(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	companyIDStr, ok := vars["company_id"]
	if !ok {
		response := map[string]string{"message": "Company ID not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userExperienceRepo := repositories.NewUserExperienceRepository()
	companyID := helper.ParseIDStringToInt(companyIDStr)
	err := userExperienceRepo.DeleteByUserIDCompanyID(userID, companyID)
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
