package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type LanguageRepository struct{}

func NewLanguageRepository() *LanguageRepository {
	return &LanguageRepository{}
}

func (repo *LanguageRepository) Create(newLanguage *models.Language) (*models.Language, error) {
	// Insert Language into database
	if err := models.DB.Create(newLanguage).Error; err != nil {
		return nil, err
	}
	return newLanguage, nil
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
	var languages []models.Language
	if err := models.DB.Find(&languages).Error; err != nil {
		return nil, err
	}
	return languages, nil
}

func (repo *LanguageRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.Language, error) {
	var languages []models.Language

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
	if err := query.Offset(offset).Limit(pageSize).Find(&languages).Error; err != nil {
		return nil, err
	}
	return languages, nil
}

func (repo *LanguageRepository) Read(id int64) (*models.Language, error) {
	var language models.Language
	if err := models.DB.First(&language, id).Error; err != nil {
		return nil, err
	}
	return &language, nil
}

func (repo *LanguageRepository) ReadByCode(code string) (*models.Language, error) {
	var language models.Language
	if err := models.DB.Where("code = ?", code).First(&language).Error; err != nil {
		return nil, err
	}
	return &language, nil
}

func (repo *LanguageRepository) ReadByName(name string) (*models.Language, error) {
	var language models.Language
	if err := models.DB.Where("name = ?", name).First(&language).Error; err != nil {
		return nil, err
	}
	return &language, nil
}

func (repo *LanguageRepository) ReadByNameOrCode(name string, code string) (*models.Language, error) {
	var language models.Language
	if err := models.DB.Where("name = ? OR code = ?", name, code).First(&language).Error; err != nil {
		return nil, err
	}
	return &language, nil
}
