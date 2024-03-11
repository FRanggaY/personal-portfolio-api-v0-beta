package public_handlers

import (
	"math"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// get public user detail project godoc
// @Summary Get public User detail project
// @Description Get Public User Detail Project with the pagination
// @Tags public-users
// @Accept json
// @Produce json
// @Param size query int false "Page size" default(5)
// @Param offset query int false "Page offset" default(1)
// @Param project_platform_id query string false "Project Platform ID"
// @Param username path string true "Username"
// @Param language_id query string true "Filter by languageId"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /public/user/{username}/project [get]
func GetPublicFilteredPaginatedUserProjectDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userNameStr, ok := vars["username"]
	if !ok {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	languageIDFilterStr := r.URL.Query().Get("language_id")
	projectPlatformIDFilterStr := r.URL.Query().Get("project_platform_id")
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	isActiveBool := helper.ParseIDStringToBool("true")

	languageIDFilter := helper.ParseIDStringToInt(languageIDFilterStr)

	var projectPlatformIDFilter *int64
	parsedID := helper.ParseIDStringToInt(projectPlatformIDFilterStr)
	if parsedID != 0 {
		projectPlatformIDFilter = &parsedID
	}

	userRepo := repositories.NewUserRepository()
	userProjectRepo := repositories.NewUserProjectRepository()

	user, err := userRepo.ReadByUsername(userNameStr)
	if err != nil {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	totalCount, err := userProjectRepo.Count(&user.ID, projectPlatformIDFilter, &isActiveBool)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	userProjects, err := userProjectRepo.ReadTranslationsByUserIDLanguageID(user.ID, projectPlatformIDFilter, languageIDFilter, &isActiveBool, pageNumber, pageSize)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch users project"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredUserProjects []models.ProjectTranslationResponse

	for _, userProject := range userProjects {

		fullImageURL := helper.GetFullImageUrl(userProject.ImageUrl, r)

		filteredUserProjects = append(filteredUserProjects, models.ProjectTranslationResponse{
			ID:                userProject.ID,
			LanguageID:        userProject.LanguageID,
			ProjectPlatformID: userProject.ProjectPlatformID,
			Name:              userProject.Name,
			Slug:              userProject.Slug,
			Description:       userProject.Description,
			ImageUrl:          fullImageURL,
			ProjectCreatedAt:  userProject.ProjectCreatedAt,
			ProjectUpdatedAt:  userProject.ProjectUpdatedAt,
			CreatedAt:         userProject.CreatedAt,
			UpdatedAt:         userProject.UpdatedAt,
		})
	}

	if filteredUserProjects == nil {
		response := map[string]string{"message": "user project not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredUserProjects,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}

// get public detail project godoc
// @Summary Get public detail project
// @Description Get Public Detail Project with the pagination
// @Tags public-users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param slug path string true "Slug"
// @Param language_id query string true "Filter by languageId"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /public/user/{username}/project/{slug} [get]
func GetPublicProjectDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userNameStr, ok := vars["username"]
	if !ok {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	slugStr, okSlug := vars["slug"]
	if !okSlug {
		response := map[string]string{"message": "Slug not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	languageIDFilterStr := r.URL.Query().Get("language_id")
	languageIDFilter := helper.ParseIDStringToInt(languageIDFilterStr)

	userRepo := repositories.NewUserRepository()
	userProjectRepo := repositories.NewUserProjectRepository()
	userProjectAttachmentRepo := repositories.NewUserProjectAttachmentRepository()

	user, err := userRepo.ReadByUsername(userNameStr)
	if err != nil {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	userProject, err := userProjectRepo.ReadByUserIDSlugLanguageID(user.ID, slugStr, languageIDFilter)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch users project"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userProjectAttachments, err := userProjectAttachmentRepo.ReadAll(&userProject.ID)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch user positions"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	var projectAttachments []map[string]interface{}

	for _, userProjectAttachment := range userProjectAttachments {
		projectAttachment := map[string]interface{}{
			"id":                    userProjectAttachment.ID,
			"title":                 userProjectAttachment.Title,
			"category":              userProjectAttachment.Category,
			"image_url":             helper.GetFullImageUrl(userProjectAttachment.ImageUrl, r),
			"url":                   userProjectAttachment.Url,
			"is_external_url":       userProjectAttachment.IsExternalUrl,
			"is_external_image_url": userProjectAttachment.IsExternalImageUrl,
		}
		projectAttachments = append(projectAttachments, projectAttachment)
	}

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"id":                  userProject.ID,
			"project_platform_id": userProject.ProjectPlatformID,
			"name":                userProject.Name,
			"slug":                userProject.Slug,
			"image_url":           helper.GetFullImageUrl(userProject.ImageUrl, r),
			"description":         userProject.Description,
			"project_created_at":  userProject.ProjectCreatedAt,
			"project_updated_at":  userProject.ProjectUpdatedAt,
			"created_at":          userProject.CreatedAt,
			"updated_at":          userProject.UpdatedAt,
			"attachments":         projectAttachments,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
