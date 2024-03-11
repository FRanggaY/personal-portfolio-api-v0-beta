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

// CreateUserEducation godoc
// @Summary Create a new user education
// @Description Create a new user education
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserEducationCreateForm true "User education input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-education [post]
func CreateUserEducation(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	// define input from json
	var userEducationInput models.UserEducationCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userEducationInput); err != nil {
		log.Fatal("Error decoding new user education: ")
	}
	defer r.Body.Close()

	userRepo := repositories.NewUserRepository()
	schoolRepo := repositories.NewSchoolRepository()
	userEducationRepo := repositories.NewUserEducationRepository()

	// validate user id
	_, user_err := userRepo.Read(userID)
	if user_err != nil {
		// Handle error
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate school id
	_, school_err := schoolRepo.Read(userEducationInput.SchoolID)
	if school_err != nil {
		// Handle error
		response := map[string]string{"message": "School ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user id and school id
	exist_data, _ := userEducationRepo.ReadByUserIDSchoolID(userID, userEducationInput.SchoolID)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Education already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newUserEducationData := models.UserEducation{
		UserID:     uint(userID),
		SchoolId:   uint(userEducationInput.SchoolID),
		MonthStart: userEducationInput.MonthStart,
		MonthEnd:   userEducationInput.MonthEnd,
		YearStart:  uint(userEducationInput.YearStart),
		YearEnd:    uint(userEducationInput.YearEnd),
	}

	// insert to database
	if newUserEducation, err := userEducationRepo.Create(&newUserEducationData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user education"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserEducation.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user Education godoc
// @Summary Delete User Education
// @Description Delete user Education
// @Tags users
// @Accept json
// @Produce json
// @Param school_id path int true "School ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-education/{school_id} [delete]
func DeleteUserEducation(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	schoolIDStr, ok := vars["school_id"]
	if !ok {
		response := map[string]string{"message": "School ID not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userEducationRepo := repositories.NewUserEducationRepository()
	schoolID := helper.ParseIDStringToInt(schoolIDStr)
	err := userEducationRepo.DeleteByUserIDSchoolID(userID, schoolID)
	if err != nil {
		response := map[string]string{"message": "User Education ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
