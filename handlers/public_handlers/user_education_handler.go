package public_handlers

import (
	"math"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// get public user detail education godoc
// @Summary Get public User detail education
// @Description Get Public User Detail education with the pagination
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
// @Router /public/user/{username}/education [get]
func GetPublicFilteredPaginatedUserEducationDetail(w http.ResponseWriter, r *http.Request) {
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

	isActiveBool := helper.ParseIDStringToBool("true")

	userRepo := repositories.NewUserRepository()
	userEducationRepo := repositories.NewUserEducationRepository()

	user, err := userRepo.ReadByUsername(userNameStr)
	if err != nil {
		response := map[string]string{"message": "Username not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	totalCount, err := userEducationRepo.Count(&user.ID, &isActiveBool)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	userEducations, err := userEducationRepo.ReadTranslationsByUserIDLanguageID(user.ID, languageIDFilter, &isActiveBool, pageNumber, pageSize)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch users education"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredUserEducations []models.EducationTranslationResponse

	for _, userEducation := range userEducations {

		fullImageURL := helper.GetFullImageUrl(userEducation.SchoolImageUrl, r)

		filteredUserEducations = append(filteredUserEducations, models.EducationTranslationResponse{
			ID:                       userEducation.ID,
			LanguageID:               userEducation.LanguageID,
			SchoolID:                 userEducation.SchoolID,
			Title:                    userEducation.Title,
			Description:              userEducation.Description,
			Category:                 userEducation.Category,
			Location:                 userEducation.Location,
			LocationType:             userEducation.LocationType,
			MonthStart:               userEducation.MonthStart,
			MonthEnd:                 userEducation.MonthEnd,
			YearStart:                userEducation.YearStart,
			YearEnd:                  userEducation.YearEnd,
			SchoolCode:               userEducation.SchoolCode,
			SchoolName:               userEducation.SchoolName,
			SchoolImageUrl:           fullImageURL,
			SchoolUrl:                userEducation.SchoolUrl,
			SchoolIsExternalUrl:      userEducation.SchoolIsExternalUrl,
			SchoolIsExternalImageUrl: userEducation.SchoolIsExternalImageUrl,
			CreatedAt:                userEducation.CreatedAt,
			UpdatedAt:                userEducation.UpdatedAt,
		})
	}

	if filteredUserEducations == nil {
		response := map[string]string{"message": "user education not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredUserEducations,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
