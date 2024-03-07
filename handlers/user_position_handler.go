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

// CreateUserPosition godoc
// @Summary Create a new user position
// @Description Create a new user position
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserPositionCreateForm true "User position input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-position [post]
func CreateUserPosition(w http.ResponseWriter, r *http.Request) {

	// define input from json
	var userPositionInput models.UserPositionCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userPositionInput); err != nil {
		log.Fatal("Error decoding new user position: ")
	}
	defer r.Body.Close()

	userRepo := repositories.NewUserRepository()
	userPositionRepo := repositories.NewUserPositionRepository()

	// validate user id
	_, user_err := userRepo.Read(userPositionInput.UserId)
	if user_err != nil {
		// Handle error
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newUserPositionData := models.UserPosition{
		UserID: uint(userPositionInput.UserId),
		Title:  userPositionInput.Title,
	}

	// insert to database
	if newUserPosition, err := userPositionRepo.Create(&newUserPositionData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user position"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserPosition.Id,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user position godoc
// @Summary Delete User position
// @Description Delete user position
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Position ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-position/{id} [delete]
func DeleteUserPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userRepoIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Position not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userPositionRepo := repositories.NewUserPositionRepository()
	userPositionID := helper.ParseIDStringToInt(userRepoIDStr)
	err := userPositionRepo.Delete(userPositionID)
	if err != nil {
		response := map[string]string{"message": "User Position ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
