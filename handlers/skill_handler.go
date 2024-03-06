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

// CreateSkill godoc
// @Summary Create a new skill
// @Description Create a new skill with file upload support
// @Tags skills
// @Accept mpfd
// @Produce json
// @Param code formData string true "Skill code"
// @Param name formData string true "Skill name"
// @Param image_file formData file true "Skill image file"
// @Param url formData string false "Skill URL"
// @Param is_external_url formData bool false "Is external URL"
// @Param is_external_image_url formData bool false "Is external image URL"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /skill [post]
func CreateSkill(w http.ResponseWriter, r *http.Request) {
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

	skillRepo := repositories.NewSkillRepository()
	// validation name and code unique
	exist_skill, _ := skillRepo.ReadByNameOrCode(name, code)
	if exist_skill != nil {
		// Handle error
		response := map[string]string{"message": "Name or Code already used"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// validation location image
	var directory = "./assets/images/skill"
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

	// Create a new skill record
	newSkill := models.Skill{
		Code:               code,
		Name:               name,
		ImageUrl:           imageUrl,
		Url:                url,
		IsExternalUrl:      isExternalUrl,
		IsExternalImageUrl: isExternalImageUrl,
	}
	// insert to database
	if newSkill, err := skillRepo.Create(&newSkill); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new skill"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newSkill.Id,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// get all skill godoc
// @Summary Get All Skill
// @Description Get All Skill with the pagination need login
// @Tags skills
// @Accept json
// @Produce json
// @Param size query int false "Page size" default(5)
// @Param offset query int false "Page offset" default(1)
// @Success 200 {object} map[string]string "Success"
// @Success 500 {object} map[string]string "Internal Server Error"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /skill [get]
func GetFilteredPaginatedSkills(w http.ResponseWriter, r *http.Request) {
	pageSize := helper.ParsePageSize(r.URL.Query().Get("size"))
	pageNumber := helper.ParsePageNumber(r.URL.Query().Get("offset"))

	skillRepo := repositories.NewSkillRepository()
	totalCount, err := skillRepo.Count()
	if err != nil {
		response := map[string]string{"message": "Failed to fetch total count"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	skills, err := skillRepo.ReadFilteredPaginated(pageSize, pageNumber)
	if err != nil {
		response := map[string]string{"message": "Failed to fetch skills"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var filteredSkills []struct {
		Id        int64     `json:"id"`
		Code      string    `json:"code"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	for _, skill := range skills {
		filteredSkills = append(filteredSkills, struct {
			Id        int64     `json:"id"`
			Code      string    `json:"code"`
			Name      string    `json:"name"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			Id:        skill.Id,
			Code:      skill.Code,
			Name:      skill.Name,
			CreatedAt: skill.CreatedAt,
			UpdatedAt: skill.UpdatedAt,
		})
	}

	if filteredSkills == nil {
		response := map[string]string{"message": "skills not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data":    filteredSkills,
		"meta": map[string]interface{}{
			"size":       pageSize,
			"offset":     pageNumber,
			"totalCount": totalCount,
			"totalPage":  totalPage,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}

// ReadSkill godoc
// @Summary Read skill
// @Description Retrieve details of a skill by its ID
// @Tags skills
// @Accept json
// @Produce json
// @Param id path int true "Skill ID"
// @Success 200 {object} models.Skill "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /skill/{id} [get]
func ReadSkill(w http.ResponseWriter, r *http.Request) {
	// Extract skill ID from the request URL
	vars := mux.Vars(r)
	skillIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Skill ID not provided", http.StatusBadRequest)
		return
	}

	skillID := helper.ParseIDStringToInt(skillIDStr)

	skillRepo := repositories.NewSkillRepository()
	skill, err := skillRepo.Read(skillID)
	if err != nil {
		response := map[string]string{"message": "Skill ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	fullImageURL := helper.GetFullImageUrl(skill.ImageUrl, r)

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"id":                    skill.Id,
			"code":                  skill.Code,
			"name":                  skill.Name,
			"is_external_url":       skill.IsExternalUrl,
			"is_external_image_url": skill.IsExternalImageUrl,
			"url":                   skill.Url,
			"created_at":            skill.CreatedAt,
			"updated_at":            skill.UpdatedAt,
			"image_url":             fullImageURL,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
