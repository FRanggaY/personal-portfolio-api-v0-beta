package public_handlers

import (
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// get public user detail godoc
// @Summary Get public User detail
// @Description Get Public User Detail with the pagination
// @Tags public-users
// @Accept json
// @Produce json
// @Param size query int false "Page size" default(5)
// @Param offset query int false "Page offset" default(1)
// @Param username path string true "Username"
// @Param language_id query string true "Filter by languageId"
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /public/user/{username} [get]
func GetPublicFilteredPaginatedUserDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userNameStr, ok := vars["username"]
	if !ok {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	languageIDFilterStr := r.URL.Query().Get("language_id")
	languageIDFilter := helper.ParseIDStringToInt(languageIDFilterStr)
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	isActiveBool := helper.ParseIDStringToBool("true")

	userRepo := repositories.NewUserRepository()
	userAttachmentRepo := repositories.NewUserAttachmentRepository()
	userPositionRepo := repositories.NewUserPositionRepository()
	userLanguageRepo := repositories.NewUserLanguageRepository()

	user, err := userRepo.ReadByUsername(userNameStr)
	if err != nil {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	userPositions, err := userPositionRepo.ReadAll(&user.ID, &isActiveBool)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch user positions"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	userAttachments, err := userAttachmentRepo.ReadAll(&user.ID, nil, &isActiveBool)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch user attachment"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userLanguages, err := userLanguageRepo.ReadTranslationsByUserIDLanguageID(user.ID, languageIDFilter, &isActiveBool, pageNumber, pageSize)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch user language"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var positions []map[string]interface{}
	var attachments []map[string]interface{}
	var languages []map[string]interface{}

	for _, userPosition := range userPositions {
		position := map[string]interface{}{
			"id":    userPosition.ID,
			"title": userPosition.Title,
		}
		positions = append(positions, position)
	}

	for _, userAttachment := range userAttachments {
		attachment := map[string]interface{}{
			"id":                    userAttachment.ID,
			"title":                 userAttachment.Title,
			"category":              userAttachment.Category,
			"image_url":             helper.GetFullImageUrl(userAttachment.ImageUrl, r),
			"url":                   userAttachment.Url,
			"is_external_url":       userAttachment.IsExternalUrl,
			"is_external_image_url": userAttachment.IsExternalImageUrl,
		}
		attachments = append(attachments, attachment)
	}

	for _, userLanguage := range userLanguages {
		language := map[string]interface{}{
			"id":          userLanguage.ID,
			"code":        userLanguage.Code,
			"title":       userLanguage.Title,
			"name":        userLanguage.Name,
			"description": userLanguage.Description,
			"logo_url":    helper.GetFullImageUrl(userLanguage.LogoUrl, r),
		}
		languages = append(languages, language)
	}

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"id":          user.ID,
			"username":    user.Username,
			"name":        user.Name,
			"created_at":  user.CreatedAt,
			"updated_at":  user.UpdatedAt,
			"positions":   positions,
			"attachments": attachments,
			"languages":   languages,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
