package handlers

import (
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// CreateLanguage godoc
// @Summary Create a new language
// @Description Create a new language with file upload support
// @Tags languages
// @Accept mpfd
// @Produce json
// @Param code formData string true "Language code"
// @Param name formData string true "Language name"
// @Param logo_url formData file true "Language logo file"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /language [post]
func CreateLanguage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Max size file only 10 mb"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	code := r.FormValue("code")
	name := r.FormValue("name")
	// Get the uploaded file
	file, header, err := r.FormFile("logo_url")
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Failed to get file"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	languageRepo := repositories.NewLanguageRepository()
	// validation name and code unique
	exist_language, _ := languageRepo.ReadByNameOrCode(name, code)
	if exist_language != nil {
		// Handle error
		response := map[string]string{"message": "Name or Code already used"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validation location image
	var directory = "./assets/images/language"
	if err := os.MkdirAll(directory, 0755); err != nil {
		// Handle error
		response := map[string]string{"message": "Failed to create directory"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	filename := header.Filename
	extension := filepath.Ext(filename)
	var finalFilename = code + extension
	imageUrl, _ := helper.UploadFile(file, directory, finalFilename)

	// Create a new language record
	newLangugage := models.Language{
		Code:    code,
		Name:    name,
		LogoUrl: imageUrl,
	}
	// insert to database
	if newLanguage, err := languageRepo.Create(&newLangugage); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new language"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newLanguage.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// get all language godoc
// @Summary Get All Language
// @Description Get All Language with the pagination need login
// @Tags languages
// @Accept json
// @Produce json
// @Param size query int false "Page size" default(5)
// @Param offset query int false "Page offset" default(1)
// @Success 200 {object} map[string][]models.LanguageResponse "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /language [get]
func GetFilteredPaginatedLanguages(w http.ResponseWriter, r *http.Request) {
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	languageRepo := repositories.NewLanguageRepository()
	totalCount, err := languageRepo.Count()
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	languages, err := languageRepo.ReadFilteredPaginated(pageSize, pageNumber)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch languages"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredLanguages []models.LanguageResponse
	for _, language := range languages {
		fullImageURL := helper.GetFullImageUrl(language.LogoUrl, r)

		filteredLanguages = append(filteredLanguages, models.LanguageResponse{
			ID:        language.ID,
			Code:      language.Code,
			Name:      language.Name,
			LogoUrl:   fullImageURL,
			CreatedAt: language.CreatedAt,
			UpdatedAt: language.UpdatedAt,
		})
	}

	if filteredLanguages == nil {
		response := map[string]string{"message": "languages not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredLanguages,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}

// ReadCompnay godoc
// @Summary Read language
// @Description Retrieve details of a language by its ID
// @Tags languages
// @Accept json
// @Produce json
// @Param id path int true "Language ID"
// @Success 200 {object} models.LanguageResponse "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /language/{id} [get]
func ReadLanguage(w http.ResponseWriter, r *http.Request) {
	// Extract language ID from the request URL
	vars := mux.Vars(r)
	languageIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Language ID not provided", http.StatusBadRequest)
		return
	}

	languageID := helper.ParseIDStringToInt(languageIDStr)

	languageRepo := repositories.NewLanguageRepository()
	language, err := languageRepo.Read(languageID)
	if err != nil {
		response := map[string]string{"message": "Language ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	fullImageURL := helper.GetFullImageUrl(language.LogoUrl, r)

	var responseData = models.LanguageResponse{
		ID:        language.ID,
		Code:      language.Code,
		Name:      language.Name,
		LogoUrl:   fullImageURL,
		CreatedAt: language.CreatedAt,
		UpdatedAt: language.UpdatedAt,
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    responseData,
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
