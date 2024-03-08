package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type SkillTranslationRepository struct{}

func NewSkillTranslationRepository() *SkillTranslationRepository {
	return &SkillTranslationRepository{}
}

func (repo *SkillTranslationRepository) Create(newData *models.SkillTranslation) (*models.SkillTranslation, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *SkillTranslationRepository) Count(languageId *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.SkillTranslation{})

	if languageId != nil {
		query = query.Where("language_id LIKE ?", languageId)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *SkillTranslationRepository) ReadAll(languageId *int64) ([]models.SkillTranslation, error) {
	query := models.DB
	var datas []models.SkillTranslation

	if languageId != nil {
		query = query.Where("language_id LIKE ?", languageId)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *SkillTranslationRepository) ReadFilteredPaginated(languageId *int64, pageSize, pageNumber int) ([]models.SkillTranslation, error) {
	var datas []models.SkillTranslation

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
	if languageId != nil {
		query = query.Where("language_id LIKE ?", languageId)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *SkillTranslationRepository) Read(id int64) (*models.SkillTranslation, error) {
	var data models.SkillTranslation
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *SkillTranslationRepository) ReadByLanguageIdSkillId(languageId int64, skillId int64) (*models.SkillTranslation, error) {
	var data models.SkillTranslation
	if err := models.DB.Where("language_id = ? AND skill_id = ?", languageId, skillId).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *SkillTranslationRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.SkillTranslation{}, id).Error; err != nil {
		return err
	}
	return nil
}
