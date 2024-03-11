package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// CreateUserProject godoc
// @Summary Create a new user project
// @Description Create a new user project with file upload support
// @Tags users
// @Accept mpfd
// @Produce json
// @Param project_platform_id formData string true "project platform id"
// @Param project_created_at formData string true "User project created at"
// @Param project_updated_at formData string true "User project updated at"
// @Param slug formData string true "User project slug"
// @Param image_file formData file true "User project image file"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-project [post]
func CreateUserProject(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)

	err := r.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Max size file only 10 mb"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	userID := jwtClaim.Id
	slug := r.FormValue("slug")
	projectPlatformIDStr := r.FormValue("project_platform_id")
	projectPlatformID := helper.ParseIDStringToInt(projectPlatformIDStr)

	projectCreatedAt := r.FormValue("project_created_at")
	projectUpdatedAt := r.FormValue("project_updated_at")
	// Get the uploaded file
	file, header, err := r.FormFile("image_file")
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Failed to get file"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	userRepo := repositories.NewUserRepository()
	projectPlatformRepo := repositories.NewProjectPlatformRepository()
	userProjectRepo := repositories.NewUserProjectRepository()
	// validation name and code unique
	_, errProjectPlatform := projectPlatformRepo.Read(projectPlatformID)
	if errProjectPlatform != nil {
		response := map[string]string{"message": "Project Platfrom ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	_, err_user := userRepo.Read(userID)
	if err_user != nil {
		// Handle error
		response := map[string]string{"message": "User not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339, projectCreatedAt)
	if err != nil {
		response := map[string]string{"message": "User not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	parsedUpdatedAt, err := time.Parse(time.RFC3339, projectUpdatedAt)
	if err != nil {
		response := map[string]string{"message": "User not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validation location image
	var directory = "./assets/images/user/project/showcase"
	if err := os.MkdirAll(directory, 0755); err != nil {
		// Handle error
		response := map[string]string{"message": "Failed to create directory"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	filename := header.Filename
	extension := filepath.Ext(filename)
	finalFilename := helper.GetStringTimeNow() + "_" + helper.ParseIDIntToString(userID) + "_" + extension
	imageUrl, _ := helper.UploadFile(file, directory, finalFilename)

	// Create a new user project record
	newUserProjectData := models.UserProject{
		ProjectPlatformID: uint(projectPlatformID),
		UserID:            uint(userID),
		Slug:              slug,
		ImageUrl:          imageUrl,
		ProjectCreatedAt:  parsedCreatedAt,
		ProjectUpdatedAt:  parsedUpdatedAt,
	}
	// insert to database
	if newUserProject, err := userProjectRepo.Create(&newUserProjectData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user project"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserProject.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user project godoc
// @Summary Delete User project
// @Description Delete user project
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Project ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Success 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-project/{id} [delete]
func DeleteUserProject(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	userRepoIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Project not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userProjectRepo := repositories.NewUserProjectRepository()
	userProjectID := helper.ParseIDStringToInt(userRepoIDStr)

	userProject, err := userProjectRepo.Read(userProjectID)
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

	errDeleteImage := helper.RemoveFile(userProject.ImageUrl)
	if errDeleteImage != nil {
		// Handle error
		response := map[string]string{"message": "Failed to delete file"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	errDelete := userProjectRepo.Delete(userProjectID)
	if errDelete != nil {
		response := map[string]string{"message": "User Project ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
