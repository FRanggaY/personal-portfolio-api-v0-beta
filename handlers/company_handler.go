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

// CreateCompany godoc
// @Summary Create a new company
// @Description Create a new company with file upload support
// @Tags companies
// @Accept mpfd
// @Produce json
// @Param code formData string true "Company code"
// @Param name formData string true "Company name"
// @Param image_file formData file true "Company image file"
// @Param url formData string false "Company URL"
// @Param is_external_url formData bool false "Is external URL"
// @Param is_external_image_url formData bool false "Is external image URL"
// @Param address formData string false "Company address"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /company [post]
func CreateCompany(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Max size file only 10 mb"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	code := r.FormValue("code")
	name := r.FormValue("name")
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

	companyRepo := repositories.NewCompanyRepository()
	// validation name and code unique
	exist_company, _ := companyRepo.ReadByNameOrCode(name, code)
	if exist_company != nil {
		// Handle error
		response := map[string]string{"message": "Name or Code already used"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validation location image
	var directory = "./assets/images/company"
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

	// Create a new company record
	newCompany := models.Company{
		Code:               code,
		Name:               name,
		ImageUrl:           imageUrl,
		Address:            address,
		IsExternalUrl:      isExternalUrl,
		IsExternalImageUrl: isExternalImageUrl,
	}
	// insert to database
	if newCompany, err := companyRepo.Create(&newCompany); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new company"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newCompany.Id,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// get all company godoc
// @Summary Get All Company
// @Description Get All Company with the pagination need login
// @Tags companies
// @Accept json
// @Produce json
// @Param size query int false "Page size" default(5)
// @Param offset query int false "Page offset" default(1)
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /company [get]
func GetFilteredPaginatedCompanies(w http.ResponseWriter, r *http.Request) {
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	companyRepo := repositories.NewCompanyRepository()
	totalCount, err := companyRepo.Count()
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	companies, err := companyRepo.ReadFilteredPaginated(pageSize, pageNumber)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch companies"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredCompanies []struct {
		Id        int64     `json:"id"`
		Code      string    `json:"code"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	for _, company := range companies {
		filteredCompanies = append(filteredCompanies, struct {
			Id        int64     `json:"id"`
			Code      string    `json:"code"`
			Name      string    `json:"name"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			Id:        company.Id,
			Code:      company.Code,
			Name:      company.Name,
			CreatedAt: company.CreatedAt,
			UpdatedAt: company.UpdatedAt,
		})
	}

	if filteredCompanies == nil {
		response := map[string]string{"message": "companies not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredCompanies,
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
// @Summary Read company
// @Description Retrieve details of a company by its ID
// @Tags companies
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} models.Company "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /company/{id} [get]
func ReadCompany(w http.ResponseWriter, r *http.Request) {
	// Extract company ID from the request URL
	vars := mux.Vars(r)
	companyIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Company ID not provided", http.StatusBadRequest)
		return
	}

	companyID := helper.ParseIDStringToInt(companyIDStr)

	companyRepo := repositories.NewCompanyRepository()
	company, err := companyRepo.Read(companyID)
	if err != nil {
		response := map[string]string{"message": "Company ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	fullImageURL := helper.GetFullImageUrl(company.ImageUrl, r)

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"id":                    company.Id,
			"code":                  company.Code,
			"name":                  company.Name,
			"is_external_url":       company.IsExternalUrl,
			"is_external_image_url": company.IsExternalImageUrl,
			"url":                   company.Url,
			"address":               company.Address,
			"created_at":            company.CreatedAt,
			"updated_at":            company.UpdatedAt,
			"image_url":             fullImageURL,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
