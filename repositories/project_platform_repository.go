package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type ProjectPlatformRepository struct{}

func NewProjectPlatformRepository() *ProjectPlatformRepository {
	return &ProjectPlatformRepository{}
}

func (repo *ProjectPlatformRepository) Create(newData *models.ProjectPlatform) (*models.ProjectPlatform, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *ProjectPlatformRepository) Count() (int, error) {
	var count int64
	query := models.DB.Model(&models.ProjectPlatform{})
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *ProjectPlatformRepository) ReadAll() ([]models.ProjectPlatform, error) {
	var datas []models.ProjectPlatform
	if err := models.DB.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *ProjectPlatformRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.ProjectPlatform, error) {
	var datas []models.ProjectPlatform

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

func (repo *ProjectPlatformRepository) Read(ID int64) (*models.ProjectPlatform, error) {
	var data models.ProjectPlatform
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *ProjectPlatformRepository) ReadByCode(code string) (*models.ProjectPlatform, error) {
	var data models.ProjectPlatform
	if err := models.DB.Where("code = ?", code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *ProjectPlatformRepository) ReadByName(name string) (*models.ProjectPlatform, error) {
	var data models.ProjectPlatform
	if err := models.DB.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *ProjectPlatformRepository) ReadByNameOrCode(name string, code string) (*models.ProjectPlatform, error) {
	var data models.ProjectPlatform
	if err := models.DB.Where("name = ? OR code = ?", name, code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
