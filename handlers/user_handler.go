package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

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
// @Router /user [get]
func GetFilteredPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	nameFilter := r.URL.Query().Get("name")
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	totalCount, err := repositories.CountUsers(nameFilter)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	users, err := repositories.ReadUsersFilteredPaginated(nameFilter, pageSize, pageNumber)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch users"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredUsers []struct {
		Username  string    `json:"username"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	for _, user := range users {
		filteredUsers = append(filteredUsers, struct {
			Username  string    `json:"username"`
			Name      string    `json:"name"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			Username:  user.Username,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
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
		response := map[string]string{"message": "User ID not found in URL"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userID := helper.ParseUserID(userIDStr)
	user, err := repositories.ReadUser(userID)
	if err != nil {
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"username":   user.Username,
			"name":       user.Name,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}

// func UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	var updatedUser models.User
// 	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
// 		http.Error(w, "Failed to decode user input", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	if err := repositories.UpdateUser(&updatedUser); err != nil {
// 		http.Error(w, "Failed to update user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("User updated successfully"))
// }

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
		response := map[string]string{"message": "User ID not found in URL"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userID := helper.ParseUserID(userIDStr)
	err := repositories.DeleteUser(userID)
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
