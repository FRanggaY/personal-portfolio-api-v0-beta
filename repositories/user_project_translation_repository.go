package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserProjectTranslationRepository struct{}

func NewUserProjectTranslationRepository() *UserProjectTranslationRepository {
	return &UserProjectTranslationRepository{}
}

func (repo *UserProjectTranslationRepository) Create(newData *models.UserProjectTranslation) (*models.UserProjectTranslation, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserProjectTranslationRepository) Count(languageID *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserProjectTranslation{})

	if languageID != nil {
		query = query.Where(helper.FilterLanguageIDEqual, languageID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserProjectTranslationRepository) ReadAll(languageID *int64) ([]models.UserProjectTranslation, error) {
	query := models.DB
	var datas []models.UserProjectTranslation

	if languageID != nil {
		query = query.Where(helper.FilterLanguageIDEqual, languageID)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserProjectTranslationRepository) ReadFilteredPaginated(languageID *int64, pageSize, pageNumber int) ([]models.UserProjectTranslation, error) {
	var datas []models.UserProjectTranslation

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

func (repo *UserProjectTranslationRepository) Read(ID int64) (*models.UserProjectTranslation, error) {
	var data models.UserProjectTranslation
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserProjectTranslationRepository) ReadByLanguageIDUserProjectID(languageID int64, userProjectID int64) (*models.UserProjectTranslation, error) {
	var data models.UserProjectTranslation
	if err := models.DB.Where("language_id = ? AND user_project_id = ?", languageID, userProjectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserProjectTranslationRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserProjectTranslation{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserProjectTranslationRepository) DeleteByLanguageIDUserProjectID(languageID int64, userProjectID int64) error {
	if err := models.DB.
		Where("language_id = ? AND user_project_id = ?", languageID, userProjectID).
		Delete(&models.UserProjectTranslation{}).
		Error; err != nil {
		return err
	}
	return nil
}
