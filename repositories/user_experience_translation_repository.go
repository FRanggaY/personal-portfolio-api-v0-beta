package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserExperienceTranslationRepository struct{}

func NewUserExperienceTranslationRepository() *UserExperienceTranslationRepository {
	return &UserExperienceTranslationRepository{}
}

func (repo *UserExperienceTranslationRepository) Create(newData *models.UserExperienceTranslation) (*models.UserExperienceTranslation, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserExperienceTranslationRepository) Count() (int, error) {
	var count int64
	query := models.DB.Model(&models.UserExperienceTranslation{})
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserExperienceTranslationRepository) ReadAll() ([]models.UserExperienceTranslation, error) {
	var datas []models.UserExperienceTranslation
	if err := models.DB.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserExperienceTranslationRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.UserExperienceTranslation, error) {
	var datas []models.UserExperienceTranslation

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

func (repo *UserExperienceTranslationRepository) Read(ID int64) (*models.UserExperienceTranslation, error) {
	var data models.UserExperienceTranslation
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserExperienceTranslationRepository) ReadByLanguageIDUserExperienceID(languageID int64, userExperienceID int64) (*models.UserExperienceTranslation, error) {
	var data models.UserExperienceTranslation
	if err := models.DB.Where("language_id = ? AND user_experience_id = ?", languageID, userExperienceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserExperienceTranslationRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserExperienceTranslation{}, ID).Error; err != nil {
		return err
	}
	return nil
}
