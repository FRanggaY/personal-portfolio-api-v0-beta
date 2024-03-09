package public_handlers

import (
	"math"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// get public user detail experience godoc
// @Summary Get public User detail experience
// @Description Get Public User Detail experience with the pagination
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
// @Router /public/user/{username}/experience [get]
func GetPublicFilteredPaginatedUserExperienceDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userNameStr, ok := vars["username"]
	if !ok {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	languageIDFilterStr := r.URL.Query().Get("language_id")
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	languageIDFilter := helper.ParseIDStringToInt(languageIDFilterStr)

	userRepo := repositories.NewUserRepository()
	userExperienceRepo := repositories.NewUserExperienceRepository()

	user, err := userRepo.ReadByUsername(userNameStr)
	if err != nil {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	totalCount, err := userExperienceRepo.Count(&user.ID)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	userExperiences, err := userExperienceRepo.ReadTranslationsByUserIDLanguageID(user.ID, languageIDFilter, pageNumber, pageSize)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch users experience"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredUserExperiences []models.ExperienceTranslationResponse

	for _, userExperience := range userExperiences {

		fullImageURL := helper.GetFullImageUrl(userExperience.CompanyImageUrl, r)

		filteredUserExperiences = append(filteredUserExperiences, models.ExperienceTranslationResponse{
			ID:                        userExperience.ID,
			LanguageID:                userExperience.LanguageID,
			CompanyID:                 userExperience.CompanyID,
			Title:                     userExperience.Title,
			Description:               userExperience.Description,
			Category:                  userExperience.Category,
			Location:                  userExperience.Location,
			LocationType:              userExperience.LocationType,
			Industry:                  userExperience.Industry,
			MonthStart:                userExperience.MonthStart,
			MonthEnd:                  userExperience.MonthEnd,
			YearStart:                 userExperience.YearStart,
			YearEnd:                   userExperience.YearEnd,
			CompanyCode:               userExperience.CompanyCode,
			CompanyName:               userExperience.CompanyName,
			CompanyImageUrl:           fullImageURL,
			CompanyUrl:                userExperience.CompanyUrl,
			CompanyIsExternalUrl:      userExperience.CompanyIsExternalUrl,
			CompanyIsExternalImageUrl: userExperience.CompanyIsExternalImageUrl,
			CreatedAt:                 userExperience.CreatedAt,
			UpdatedAt:                 userExperience.UpdatedAt,
		})
	}

	if filteredUserExperiences == nil {
		response := map[string]string{"message": "user experience not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredUserExperiences,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
