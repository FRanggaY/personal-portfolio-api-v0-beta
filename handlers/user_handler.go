package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

var messageUserIdDetailNotFound = "User ID not found in URL"

// get all user godoc
// @Summary Get All User
// @Description Get All User with the pagination need login
// @Tags users
// @Accept json
// @Produce json
// @Param size query int false "Page size" default(5)
// @Param offset query int false "Page offset" default(1)
// @Param name query string false "Filter by name"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user [get]
func GetFilteredPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	nameFilter := r.URL.Query().Get("name")
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	userRepo := repositories.NewUserRepository()
	totalCount, err := userRepo.Count(nameFilter)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	users, err := userRepo.ReadFilteredPaginated(nameFilter, pageSize, pageNumber)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch users"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredUsers []struct {
		Id        int64     `json:"id"`
		Username  string    `json:"username"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	for _, user := range users {
		filteredUsers = append(filteredUsers, struct {
			Id        int64     `json:"id"`
			Username  string    `json:"username"`
			Name      string    `json:"name"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			Id:        user.Id,
			Username:  user.Username,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	if filteredUsers == nil {
		response := map[string]string{"message": "user not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredUsers,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}

// get user godoc
// @Summary Get User
// @Description Get user with provided detail by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": messageUserIdDetailNotFound}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userRepo := repositories.NewUserRepository()
	userID := helper.ParseIDStringToInt(userIDStr)
	user, err := userRepo.Read(userID)
	if err != nil {
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"id":         user.Id,
			"username":   user.Username,
			"name":       user.Name,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}

// update user godoc
// @Summary Update User
// @Description Update user with provided detail by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body models.UserEditForm true "User input"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": messageUserIdDetailNotFound}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userID := helper.ParseIDStringToInt(userIDStr)

	var updatedUser models.UserEditForm
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		response := map[string]string{"message": "Failed to decode user input"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	userRepo := repositories.NewUserRepository()
	// validate username unique in other user
	exist_user, _ := userRepo.ReadByUsername(updatedUser.Username)
	if exist_user != nil && exist_user.Id != userID {
		// Handle error
		response := map[string]string{"message": "Username already used"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if err := userRepo.Update(userID, &updatedUser); err != nil {
		response := map[string]string{"message": "Failed to update user"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

// delete user godoc
// @Summary Delete User
// @Description Delete user with provided detail by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": messageUserIdDetailNotFound}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userRepo := repositories.NewUserRepository()
	userID := helper.ParseIDStringToInt(userIDStr)
	err := userRepo.Delete(userID)
	if err != nil {
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

// CreateUserSkill godoc
// @Summary Create a new user skill
// @Description Create a new user skill
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.UserSkillCreateForm true "User input"
// @Param user_id formData string true "User ID"
// @Param skill_id formData string true "Skill ID"
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
	_, user_err := userRepo.Read(userSkillInput.UserId)
	if user_err != nil {
		// Handle error
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate skill id
	_, skill_err := skillRepo.Read(userSkillInput.SkillId)
	if skill_err != nil {
		// Handle error
		response := map[string]string{"message": "Skill ID not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validate user id and skill id
	exist_data, _ := userSkillRepo.ReadByUserIdSkillId(userSkillInput.UserId, userSkillInput.SkillId)
	if exist_data != nil {
		// Handle error
		response := map[string]string{"message": "Skill already added"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	newUserSkillData := models.UserSkill{
		UserID:  uint(userSkillInput.UserId),
		SkillId: uint(userSkillInput.SkillId),
	}

	// insert to database
	if newUserSkill, err := userSkillRepo.Create(&newUserSkillData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserSkill.Id,
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
	userIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": messageUserIdDetailNotFound}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userSkillRepo := repositories.NewUserSkillRepository()
	userID := helper.ParseIDStringToInt(userIDStr)
	err := userSkillRepo.Delete(userID)
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
