package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserLanguageRepository struct{}

func NewUserLanguageRepository() *UserLanguageRepository {
	return &UserLanguageRepository{}
}

func (repo *UserLanguageRepository) Create(newData *models.UserLanguage) (*models.UserLanguage, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserLanguageRepository) Count(userID *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserLanguage{})

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserLanguageRepository) ReadAll(userID *int64) ([]models.UserLanguage, error) {
	query := models.DB
	var datas []models.UserLanguage

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserLanguageRepository) ReadFilteredPaginated(userID *int64, pageSize, pageNumber int) ([]models.UserLanguage, error) {
	var datas []models.UserLanguage

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
	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserLanguageRepository) Read(ID int64) (*models.UserLanguage, error) {
	var data models.UserLanguage
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserLanguageRepository) ReadByUserIDLanguageID(userID int64, languageID int64) (*models.UserLanguage, error) {
	var data models.UserLanguage
	if err := models.DB.Where("user_id = ? AND language_id = ?", userID, languageID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserLanguageRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserLanguage{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserLanguageRepository) DeleteByUserIDLanguageID(userID int64, languageID int64) error {
	if err := models.DB.
		Where("user_id = ? AND language_id = ?", userID, languageID).
		Delete(&models.UserLanguage{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserLanguageRepository) ReadTranslationsByUserIDLanguageID(userID int64, languageID int64, pageNumber int, pageSize int) ([]models.UserLanguageTranslationResponse, error) {
	var languages []models.UserLanguageTranslationResponse

	// default
	if pageSize <= 0 {
		pageSize = 5
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	// calculate offset
	offset := (pageNumber - 1) * pageSize

	query := models.DB.
		Table("user_language_translations").
		Select(`
            languages.*, 
            IFNULL(user_language_translations.title, '') AS title, 
            IFNULL(user_language_translations.description, '') AS description, 
            IFNULL(user_language_translations.language_id, '') AS language_id
        `).
		Joins("LEFT JOIN user_languages ON user_language_translations.user_language_id = user_languages.id").
		Joins("LEFT JOIN languages ON user_languages.language_id = languages.id").
		Where("user_languages.user_id = ?", userID).
		Where("user_language_translations.language_id = ?", languageID).
		Limit(pageSize).Offset(offset)

	if err := query.Find(&languages).Error; err != nil {
		return nil, err
	}

	return languages, nil
}
