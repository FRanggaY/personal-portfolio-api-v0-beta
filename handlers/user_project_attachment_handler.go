package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// CreateUserProjectAttachment godoc
// @Summary Create a new user project attachment
// @Description Create a new user project attachment with file upload support
// @Tags users
// @Accept mpfd
// @Produce json
// @Param user_project_id formData string true "user project id"
// @Param title formData string true "User attachment title"
// @Param category formData string true "User attachment category"
// @Param image_file formData file true "User attachment image file"
// @Param url formData string false "User attachment URL"
// @Param is_external_url formData bool false "Is external URL"
// @Param is_external_image_url formData bool false "Is external image URL"
// @Success 201 {object} map[string]string "Created"
// @Success 403 {object} map[string]string "Forbidden"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-project-attachment [post]
func CreateUserProjectAttachment(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)

	err := r.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Max size file only 10 mb"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	userID := jwtClaim.Id
	title := r.FormValue("title")
	category := r.FormValue("category")
	url := r.FormValue("url")
	isExternalUrl := helper.ParseIDStringToBool(r.FormValue("is_external_url"))
	isExternalImageUrl := helper.ParseIDStringToBool(r.FormValue("is_external_image_url"))
	userProjectIDStr := r.FormValue("user_project_id")
	userProjectID := helper.ParseIDStringToInt(userProjectIDStr)

	// Get the uploaded file
	file, header, err := r.FormFile("image_file")
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Failed to get file"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	userProjectRepo := repositories.NewUserProjectRepository()
	userProjectAttachmentRepo := repositories.NewUserProjectAttachmentRepository()
	// validation name and code unique
	userProject, errUserProject := userProjectRepo.Read(userProjectID)
	if errUserProject != nil {
		response := map[string]string{"message": "User Project ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	if userProject.UserID != uint(userID) {
		response := map[string]string{"message": "Not allowing create"}
		helper.ResponseJSON(w, http.StatusForbidden, response)
		return
	}

	// validation location image
	var directory = "./assets/images/user/project/showcase_attachment"
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
	newUserProjectAttachmentData := models.UserProjectAttachment{
		UserProjectID:      uint(userProjectID),
		Title:              title,
		Category:           category,
		ImageUrl:           imageUrl,
		Url:                url,
		IsExternalUrl:      isExternalUrl,
		IsExternalImageUrl: isExternalImageUrl,
	}
	// insert to database
	if newUserProjectAttachment, err := userProjectAttachmentRepo.Create(&newUserProjectAttachmentData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user project attachment"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserProjectAttachment.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user project attachment godoc
// @Summary Delete User project attachment
// @Description Delete user project attachment
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Project Attachment ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Success 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-project-attachment/{id} [delete]
func DeleteUserProjectAttachment(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	userProjectAttachmentIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Project Attachment not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userProjectRepo := repositories.NewUserProjectRepository()
	userProjectAttachmentRepo := repositories.NewUserProjectAttachmentRepository()
	userProjectAttachmentID := helper.ParseIDStringToInt(userProjectAttachmentIDStr)

	userProjectAttachment, err := userProjectAttachmentRepo.Read(userProjectAttachmentID)
	if err != nil {
		response := map[string]string{"message": "User Project Attachment ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	userProject, err := userProjectRepo.Read(int64(userProjectAttachment.UserProjectID))
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

	errDeleteImage := helper.RemoveFile(userProjectAttachment.ImageUrl)
	if errDeleteImage != nil {
		// Handle error
		response := map[string]string{"message": "Failed to delete file"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	errDelete := userProjectAttachmentRepo.Delete(userProjectAttachmentID)
	if errDelete != nil {
		response := map[string]string{"message": "User Project Attachment ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
