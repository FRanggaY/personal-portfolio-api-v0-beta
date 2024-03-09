package public_handlers

import (
	"math"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// get public user detail skill godoc
// @Summary Get public User detail skill
// @Description Get Public User Detail Skill with the pagination
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
// @Router /public/user/{username}/skill [get]
func GetPublicFilteredPaginatedUserSkillDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userNameStr, ok := vars["username"]
	if !ok {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	languageIdFilterStr := r.URL.Query().Get("language_id")
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	languageIDFilter := helper.ParseIDStringToInt(languageIdFilterStr)

	userRepo := repositories.NewUserRepository()
	userSkillRepo := repositories.NewUserSkillRepository()

	user, err := userRepo.ReadByUsername(userNameStr)
	if err != nil {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	totalCount, err := userSkillRepo.Count(&user.Id)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	userSkills, err := userSkillRepo.ReadTranslationsByUserIDLanguageID(user.Id, languageIDFilter, pageNumber, pageSize)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch users skill"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredUserSkills []models.SkillTranslationResponse

	for _, userSkill := range userSkills {

		fullImageURL := helper.GetFullImageUrl(userSkill.ImageUrl, r)

		filteredUserSkills = append(filteredUserSkills, models.SkillTranslationResponse{
			ID:                 userSkill.ID,
			LanguageID:         userSkill.LanguageID,
			Code:               userSkill.Code,
			Name:               userSkill.Name,
			ImageUrl:           fullImageURL,
			Url:                userSkill.Url,
			IsExternalUrl:      userSkill.IsExternalUrl,
			IsExternalImageUrl: userSkill.IsExternalImageUrl,
			Description:        userSkill.Description,
			CreatedAt:          userSkill.CreatedAt,
			UpdatedAt:          userSkill.UpdatedAt,
		})
	}

	if filteredUserSkills == nil {
		response := map[string]string{"message": "user skill not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredUserSkills,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
