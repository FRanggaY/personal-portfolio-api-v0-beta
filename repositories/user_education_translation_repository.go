package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserEducationTranslationRepository struct{}

func NewUserEducationTranslationRepository() *UserEducationTranslationRepository {
	return &UserEducationTranslationRepository{}
}

func (repo *UserEducationTranslationRepository) Create(newData *models.UserEducationTranslation) (*models.UserEducationTranslation, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserEducationTranslationRepository) Count() (int, error) {
	var count int64
	query := models.DB.Model(&models.UserEducationTranslation{})
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserEducationTranslationRepository) ReadAll() ([]models.UserEducationTranslation, error) {
	var datas []models.UserEducationTranslation
	if err := models.DB.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserEducationTranslationRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.UserEducationTranslation, error) {
	var datas []models.UserEducationTranslation

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

func (repo *UserEducationTranslationRepository) Read(id int64) (*models.UserEducationTranslation, error) {
	var data models.UserEducationTranslation
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserEducationTranslationRepository) ReadByLanguageIdUserEducationId(languageId int64, userEducationId int64) (*models.UserEducationTranslation, error) {
	var data models.UserEducationTranslation
	if err := models.DB.Where("language_id = ? AND user_education_id = ?", languageId, userEducationId).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserEducationTranslationRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserEducationTranslation{}, id).Error; err != nil {
		return err
	}
	return nil
}
