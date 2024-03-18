package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type ProjectPlatformTranslationRepository struct{}

func NewProjectPlatformTranslationRepository() *ProjectPlatformTranslationRepository {
	return &ProjectPlatformTranslationRepository{}
}

func (repo *ProjectPlatformTranslationRepository) Create(newData *models.ProjectPlatformTranslation) (*models.ProjectPlatformTranslation, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *ProjectPlatformTranslationRepository) Count(languageID *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.ProjectPlatformTranslation{})

	if languageID != nil {
		query = query.Where(helper.FilterLanguageIDEqual, languageID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *ProjectPlatformTranslationRepository) ReadAll(languageID *int64) ([]models.ProjectPlatformTranslation, error) {
	query := models.DB
	var datas []models.ProjectPlatformTranslation

	if languageID != nil {
		query = query.Where(helper.FilterLanguageIDEqual, languageID)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *ProjectPlatformTranslationRepository) ReadFilteredPaginated(languageID *int64, pageSize, pageNumber int) ([]models.ProjectPlatformTranslation, error) {
	var datas []models.ProjectPlatformTranslation

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
	if languageID != nil {
		query = query.Where(helper.FilterLanguageIDEqual, languageID)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *ProjectPlatformTranslationRepository) Read(ID int64) (*models.ProjectPlatformTranslation, error) {
	var data models.ProjectPlatformTranslation
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *ProjectPlatformTranslationRepository) ReadByLanguageIDProjectPlatformID(languageID int64, projectPlatformID int64) (*models.ProjectPlatformTranslation, error) {
	var data models.ProjectPlatformTranslation
	if err := models.DB.Where("language_id = ? AND project_platform_id = ?", languageID, projectPlatformID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *ProjectPlatformTranslationRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.ProjectPlatformTranslation{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProjectPlatformTranslationRepository) DeleteByLanguageIDProjectPlatformID(languageID int64, projectPlatformID int64) error {
	if err := models.DB.
		Where("language_id = ? AND project_platform_id = ?", languageID, projectPlatformID).
		Delete(&models.ProjectPlatformTranslation{}).
		Error; err != nil {
		return err
	}
	return nil
}
