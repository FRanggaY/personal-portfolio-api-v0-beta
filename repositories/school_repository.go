package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type SchoolRepository struct{}

func NewSchoolRepository() *SchoolRepository {
	return &SchoolRepository{}
}

func (repo *SchoolRepository) Create(newSchool *models.School) (*models.School, error) {
	// Insert school into database
	if err := models.DB.Create(newSchool).Error; err != nil {
		return nil, err
	}
	return newSchool, nil
}

func (repo *SchoolRepository) Count() (int, error) {
	var count int64
	query := models.DB.Model(&models.School{})
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *SchoolRepository) ReadAll() ([]models.School, error) {
	var schools []models.School
	if err := models.DB.Find(&schools).Error; err != nil {
		return nil, err
	}
	return schools, nil
}

func (repo *SchoolRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.School, error) {
	var schools []models.School

	// default
	if pageSize <= 0 {
		pageSize = 5
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	// calculate offset
	offset := (pageNumber - 1) * pageSize

	query := models.DB

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&schools).Error; err != nil {
		return nil, err
	}
	return schools, nil
}

func (repo *SchoolRepository) Read(id int64) (*models.School, error) {
	var school models.School
	if err := models.DB.First(&school, id).Error; err != nil {
		return nil, err
	}
	return &school, nil
}
