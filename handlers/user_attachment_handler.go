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

// CreateUserAttachment godoc
// @Summary Create a new user attachment
// @Description Create a new user attachment with file upload support
// @Tags users
// @Accept mpfd
// @Produce json
// @Param title formData string true "User attachment title"
// @Param category formData string true "User attachment category"
// @Param image_file formData file true "User attachment image file"
// @Param url formData string false "User attachment URL"
// @Param is_external_url formData bool false "Is external URL"
// @Param is_external_image_url formData bool false "Is external image URL"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /user-attachment [post]
func CreateUserAttachment(w http.ResponseWriter, r *http.Request) {
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
	userAttachmentRepo := repositories.NewUserAttachmentRepository()
	// validation name and code unique
	_, err_user := userRepo.Read(userID)
	if err_user != nil {
		// Handle error
		response := map[string]string{"message": "User not found"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validation location image
	var directory = "./assets/images/user/attachment"
	if err := os.MkdirAll(directory, 0755); err != nil {
		// Handle error
		response := map[string]string{"message": "Failed to create directory"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	filename := header.Filename
	extension := filepath.Ext(filename)
	finalFilename := helper.GetStringTimeNow() + "_" + helper.ParseIDIntToString(userID) + "_" + title + extension
	imageUrl, _ := helper.UploadFile(file, directory, finalFilename)

	// Create a new user attachment record
	newUserAttachmentData := models.UserAttachment{
		UserID:             uint(userID),
		Title:              title,
		Category:           category,
		ImageUrl:           imageUrl,
		Url:                url,
		IsExternalUrl:      isExternalUrl,
		IsExternalImageUrl: isExternalImageUrl,
	}
	// insert to database
	if newUserAttachment, err := userAttachmentRepo.Create(&newUserAttachmentData); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user attachment"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUserAttachment.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// delete user attachment godoc
// @Summary Delete User attachment
// @Description Delete user attachment
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User Attachment ID"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Success 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /user-attachment/{id} [delete]
func DeleteUserAttachment(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)
	userID := jwtClaim.Id

	vars := mux.Vars(r)
	userRepoIDStr, ok := vars["id"]
	if !ok {
		response := map[string]string{"message": "User Attachment not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userAttachmentRepo := repositories.NewUserAttachmentRepository()
	userAttachmentID := helper.ParseIDStringToInt(userRepoIDStr)

	userAttachment, err := userAttachmentRepo.Read(userAttachmentID)
	if err != nil {
		response := map[string]string{"message": "User Attachment ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	if userAttachment.UserID != uint(userID) {
		response := map[string]string{"message": "Not allowing delete"}
		helper.ResponseJSON(w, http.StatusForbidden, response)
		return
	}

	errDeleteImage := helper.RemoveFile(userAttachment.ImageUrl)
	if errDeleteImage != nil {
		// Handle error
		response := map[string]string{"message": "Failed to delete file"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	errDelete := userAttachmentRepo.Delete(userAttachmentID)
	if errDelete != nil {
		response := map[string]string{"message": "User Attachment ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
