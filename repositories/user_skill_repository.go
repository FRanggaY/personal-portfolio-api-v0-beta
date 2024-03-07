package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserSkillRepository struct{}

func NewUserSkillRepository() *UserSkillRepository {
	return &UserSkillRepository{}
}

func (repo *UserSkillRepository) Create(newUserSkill *models.UserSkill) (*models.UserSkill, error) {
	// Insert skill into database
	if err := models.DB.Create(newUserSkill).Error; err != nil {
		return nil, err
	}
	return newUserSkill, nil
}

func (repo *UserSkillRepository) Count(userId *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserSkill{})

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserSkillRepository) ReadAll(userId *int64) ([]models.UserSkill, error) {
	query := models.DB
	var userSkills []models.UserSkill

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Find(&userSkills).Error; err != nil {
		return nil, err
	}
	return userSkills, nil
}

func (repo *UserSkillRepository) ReadFilteredPaginated(userId *int64, pageSize, pageNumber int) ([]models.UserSkill, error) {
	var userSkills []models.UserSkill

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
	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&userSkills).Error; err != nil {
		return nil, err
	}
	return userSkills, nil
}

func (repo *UserSkillRepository) Read(id int64) (*models.UserSkill, error) {
	var userSkill models.UserSkill
	if err := models.DB.First(&userSkill, id).Error; err != nil {
		return nil, err
	}
	return &userSkill, nil
}

func (repo *UserSkillRepository) ReadByUserIdSkillId(userId int64, skillId int64) (*models.UserSkill, error) {
	var userSkill models.UserSkill
	if err := models.DB.Where("user_id = ? AND skill_id = ?", userId, skillId).First(&userSkill).Error; err != nil {
		return nil, err
	}
	return &userSkill, nil
}

func (repo *UserSkillRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserSkill{}, id).Error; err != nil {
		return err
	}
	return nil
}
