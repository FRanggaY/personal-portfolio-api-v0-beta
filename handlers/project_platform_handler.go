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

// CreateProjectPlatform godoc
// @Summary Create a new project platform
// @Description Create a new projectPlatform with file upload support
// @Tags project-platforms
// @Accept mpfd
// @Produce json
// @Param code formData string true "Project Platform code"
// @Param name formData string true "Project Platform name"
// @Param image_file formData file true "Project Platform image file"
// @Param url formData string false "Project Platform URL"
// @Param is_external_url formData bool false "Is external URL"
// @Param is_external_image_url formData bool false "Is external image URL"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /project-platform [post]
func CreateProjectPlatform(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Max size file only 10 mb"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	code := r.FormValue("code")
	name := r.FormValue("name")
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

	projectPlatformRepo := repositories.NewProjectPlatformRepository()
	// validation name and code unique
	exist_projectPlatform, _ := projectPlatformRepo.ReadByNameOrCode(name, code)
	if exist_projectPlatform != nil {
		// Handle error
		response := map[string]string{"message": "Name or Code already used"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validation location image
	var directory = "./assets/images/project/platform"
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

	// Create a new projectPlatform record
	newProjectPlatform := models.ProjectPlatform{
		Code:               code,
		Name:               name,
		ImageUrl:           imageUrl,
		Url:                url,
		IsExternalUrl:      isExternalUrl,
		IsExternalImageUrl: isExternalImageUrl,
	}
	// insert to database
	if newProjectPlatform, err := projectPlatformRepo.Create(&newProjectPlatform); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new projectPlatform"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newProjectPlatform.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// get all projectPlatform godoc
// @Summary Get All Project Platform
// @Description Get All ProjectPlatform with the pagination need login
// @Tags project-platforms
// @Accept json
// @Produce json
// @Param size query int false "Page size" default(5)
// @Param offset query int false "Page offset" default(1)
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /project-platform [get]
func GetFilteredPaginatedProjectPlatforms(w http.ResponseWriter, r *http.Request) {
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	projectPlatformRepo := repositories.NewProjectPlatformRepository()
	totalCount, err := projectPlatformRepo.Count()
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	projectPlatforms, err := projectPlatformRepo.ReadFilteredPaginated(pageSize, pageNumber)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch projectPlatforms"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredProjectPlatforms []models.ProjectPlatformAllResponse
	for _, projectPlatform := range projectPlatforms {
		fullImageURL := helper.GetFullImageUrl(projectPlatform.ImageUrl, r)

		filteredProjectPlatforms = append(filteredProjectPlatforms, models.ProjectPlatformAllResponse{
			ID:        projectPlatform.ID,
			Code:      projectPlatform.Code,
			Name:      projectPlatform.Name,
			ImageUrl:  fullImageURL,
			CreatedAt: projectPlatform.CreatedAt,
			UpdatedAt: projectPlatform.UpdatedAt,
		})
	}

	if filteredProjectPlatforms == nil {
		response := map[string]string{"message": "projectPlatforms not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredProjectPlatforms,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}

// ReadProjectPlatform godoc
// @Summary Read project Platform
// @Description Retrieve details of a project platform by its ID
// @Tags project-platforms
// @Accept json
// @Produce json
// @Param id path int true "Project Platform ID"
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /project-platform/{id} [get]
func ReadProjectPlatform(w http.ResponseWriter, r *http.Request) {
	// Extract projectPlatform ID from the request URL
	vars := mux.Vars(r)
	projectPlatformIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ProjectPlatform ID not provided", http.StatusBadRequest)
		return
	}

	projectPlatformID := helper.ParseIDStringToInt(projectPlatformIDStr)

	projectPlatformRepo := repositories.NewProjectPlatformRepository()
	projectPlatform, err := projectPlatformRepo.Read(projectPlatformID)
	if err != nil {
		response := map[string]string{"message": "ProjectPlatform ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	fullImageURL := helper.GetFullImageUrl(projectPlatform.ImageUrl, r)

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"id":                    projectPlatform.ID,
			"code":                  projectPlatform.Code,
			"name":                  projectPlatform.Name,
			"is_external_url":       projectPlatform.IsExternalUrl,
			"is_external_image_url": projectPlatform.IsExternalImageUrl,
			"url":                   projectPlatform.Url,
			"created_at":            projectPlatform.CreatedAt,
			"updated_at":            projectPlatform.UpdatedAt,
			"image_url":             fullImageURL,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
