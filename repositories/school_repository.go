package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type SchoolRepository struct{}

func NewSchoolRepository() *SchoolRepository {
	return &SchoolRepository{}
}

func (repo *SchoolRepository) Create(newData *models.School) (*models.School, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
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
	var datas []models.School
	if err := models.DB.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *SchoolRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.School, error) {
	var datas []models.School

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
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *SchoolRepository) Read(id int64) (*models.School, error) {
	var data models.School
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *SchoolRepository) ReadByCode(code string) (*models.School, error) {
	var data models.School
	if err := models.DB.Where("code = ?", code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *SchoolRepository) ReadByName(name string) (*models.School, error) {
	var data models.School
	if err := models.DB.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *SchoolRepository) ReadByNameOrCode(name string, code string) (*models.School, error) {
	var data models.School
	if err := models.DB.Where("name = ? OR code = ?", name, code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
