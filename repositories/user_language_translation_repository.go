package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserLanguageTranslationRepository struct{}

func NewUserLanguageTranslationRepository() *UserLanguageTranslationRepository {
	return &UserLanguageTranslationRepository{}
}

func (repo *UserLanguageTranslationRepository) Create(newData *models.UserLanguageTranslation) (*models.UserLanguageTranslation, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserLanguageTranslationRepository) Count(languageID *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserLanguageTranslation{})

	if languageID != nil {
		query = query.Where(helper.FilterLanguageIDEqual, languageID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserLanguageTranslationRepository) ReadAll(languageID *int64) ([]models.UserLanguageTranslation, error) {
	query := models.DB
	var datas []models.UserLanguageTranslation

	if languageID != nil {
		query = query.Where(helper.FilterLanguageIDEqual, languageID)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserLanguageTranslationRepository) ReadFilteredPaginated(languageID *int64, pageSize, pageNumber int) ([]models.UserLanguageTranslation, error) {
	var datas []models.UserLanguageTranslation

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

func (repo *UserLanguageTranslationRepository) Read(ID int64) (*models.UserLanguageTranslation, error) {
	var data models.UserLanguageTranslation
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserLanguageTranslationRepository) ReadByLanguageIDUserLanguageID(languageID int64, userLanguageID int64) (*models.UserLanguageTranslation, error) {
	var data models.UserLanguageTranslation
	if err := models.DB.Where("language_id = ? AND user_language_id = ?", languageID, userLanguageID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserLanguageTranslationRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserLanguageTranslation{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserLanguageTranslationRepository) DeleteByLanguageIDUserLanguageID(languageID int64, userLanguageID int64) error {
	if err := models.DB.
		Where("language_id = ? AND user_language_id = ?", languageID, userLanguageID).
		Delete(&models.UserLanguageTranslation{}).
		Error; err != nil {
		return err
	}
	return nil
}
