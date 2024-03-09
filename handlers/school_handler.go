package handlers

import (
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/gorilla/mux"
)

// CreateSchool godoc
// @Summary Create a new school
// @Description Create a new school with file upload support
// @Tags schools
// @Accept mpfd
// @Produce json
// @Param code formData string true "School code"
// @Param name formData string true "School name"
// @Param image_file formData file true "School image file"
// @Param url formData string false "School URL"
// @Param is_external_url formData bool false "Is external URL"
// @Param is_external_image_url formData bool false "Is external image URL"
// @Param address formData string false "School address"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /school [post]
func CreateSchool(w http.ResponseWriter, r *http.Request) {
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
	address := r.FormValue("address")
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

	schoolRepo := repositories.NewSchoolRepository()
	// validation name and code unique
	exist_school, _ := schoolRepo.ReadByNameOrCode(name, code)
	if exist_school != nil {
		// Handle error
		response := map[string]string{"message": "Name or Code already used"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validation location image
	var directory = "./assets/images/school"
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

	// Create a new school record
	newSchool := models.School{
		Code:               code,
		Name:               name,
		ImageUrl:           imageUrl,
		Url:                url,
		Address:            address,
		IsExternalUrl:      isExternalUrl,
		IsExternalImageUrl: isExternalImageUrl,
	}
	// insert to database
	if newSchool, err := schoolRepo.Create(&newSchool); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new school"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newSchool.ID,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// get all school godoc
// @Summary Get All School
// @Description Get All School with the pagination need login
// @Tags schools
// @Accept json
// @Produce json
// @Param size query int false "Page size" default(5)
// @Param offset query int false "Page offset" default(1)
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /school [get]
func GetFilteredPaginatedSchools(w http.ResponseWriter, r *http.Request) {
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	schoolRepo := repositories.NewSchoolRepository()
	totalCount, err := schoolRepo.Count()
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	schools, err := schoolRepo.ReadFilteredPaginated(pageSize, pageNumber)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch schools"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredSchools []struct {
		ID        int64     `json:"id"`
		Code      string    `json:"code"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	for _, school := range schools {
		filteredSchools = append(filteredSchools, struct {
			ID        int64     `json:"id"`
			Code      string    `json:"code"`
			Name      string    `json:"name"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			ID:        school.ID,
			Code:      school.Code,
			Name:      school.Name,
			CreatedAt: school.CreatedAt,
			UpdatedAt: school.UpdatedAt,
		})
	}

	if filteredSchools == nil {
		response := map[string]string{"message": "schools not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredSchools,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}

// ReadSchool godoc
// @Summary Read school
// @Description Retrieve details of a school by its ID
// @Tags schools
// @Accept json
// @Produce json
// @Param id path int true "School ID"
// @Success 200 {object} models.School "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /school/{id} [get]
func ReadSchool(w http.ResponseWriter, r *http.Request) {
	// Extract school ID from the request URL
	vars := mux.Vars(r)
	schoolIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "School ID not provided", http.StatusBadRequest)
		return
	}

	schoolID := helper.ParseIDStringToInt(schoolIDStr)

	schoolRepo := repositories.NewSchoolRepository()
	school, err := schoolRepo.Read(schoolID)
	if err != nil {
		response := map[string]string{"message": "School ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	fullImageURL := helper.GetFullImageUrl(school.ImageUrl, r)

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"id":                    school.ID,
			"code":                  school.Code,
			"name":                  school.Name,
			"is_external_url":       school.IsExternalUrl,
			"is_external_image_url": school.IsExternalImageUrl,
			"url":                   school.Url,
			"address":               school.Address,
			"created_at":            school.CreatedAt,
			"updated_at":            school.UpdatedAt,
			"image_url":             fullImageURL,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
