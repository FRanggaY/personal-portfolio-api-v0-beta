package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type LanguageRepository struct{}

func NewLanguageRepository() *LanguageRepository {
	return &LanguageRepository{}
}

func (repo *LanguageRepository) Create(newData *models.Language) (*models.Language, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *LanguageRepository) Count() (int, error) {
	var count int64
	query := models.DB.Model(&models.Language{})
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *LanguageRepository) ReadAll() ([]models.Language, error) {
	var datas []models.Language
	if err := models.DB.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *LanguageRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.Language, error) {
	var datas []models.Language

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

func (repo *LanguageRepository) Read(id int64) (*models.Language, error) {
	var data models.Language
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *LanguageRepository) ReadByCode(code string) (*models.Language, error) {
	var data models.Language
	if err := models.DB.Where("code = ?", code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *LanguageRepository) ReadByName(name string) (*models.Language, error) {
	var data models.Language
	if err := models.DB.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *LanguageRepository) ReadByNameOrCode(name string, code string) (*models.Language, error) {
	var data models.Language
	if err := models.DB.Where("name = ? OR code = ?", name, code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
