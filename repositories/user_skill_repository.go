package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserSkillRepository struct{}

func NewUserSkillRepository() *UserSkillRepository {
	return &UserSkillRepository{}
}

func (repo *UserSkillRepository) Create(newData *models.UserSkill) (*models.UserSkill, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserSkillRepository) Count(userID *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserSkill{})

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserSkillRepository) ReadAll(userID *int64) ([]models.UserSkill, error) {
	query := models.DB
	var datas []models.UserSkill

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserSkillRepository) ReadFilteredPaginated(userID *int64, pageSize, pageNumber int) ([]models.UserSkill, error) {
	var datas []models.UserSkill

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

func (repo *UserSkillRepository) Read(ID int64) (*models.UserSkill, error) {
	var data models.UserSkill
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserSkillRepository) ReadByUserIDSkillID(userID int64, skillID int64) (*models.UserSkill, error) {
	var data models.UserSkill
	if err := models.DB.Where("user_id = ? AND skill_id = ?", userID, skillID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserSkillRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserSkill{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserSkillRepository) ReadTranslationsByUserIDLanguageID(userID int64, languageID int64, pageNumber int, pageSize int) ([]models.SkillTranslationResponse, error) {
	var skills []models.SkillTranslationResponse

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
		Table("skill_translations").
		Select(`
            skills.*, 
            IFNULL(skill_translations.description, '') AS description, 
            IFNULL(skill_translations.language_id, '') AS language_id
        `).
		Joins("LEFT JOIN skills ON skill_translations.skill_id = skills.id").
		Joins("LEFT JOIN user_skills ON skills.id = user_skills.skill_id").
		Where("user_skills.user_id = ?", userID).
		Where("skill_translations.language_id = ?", languageID).
		Limit(pageSize).Offset(offset)

	if err := query.Find(&skills).Error; err != nil {
		return nil, err
	}

	return skills, nil
}
