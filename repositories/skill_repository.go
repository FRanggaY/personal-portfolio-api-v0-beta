package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type SkillRepository struct{}

func NewSkillRepository() *SkillRepository {
	return &SkillRepository{}
}

func (repo *SkillRepository) Create(newSkill *models.Skill) (*models.Skill, error) {
	// Insert skill into database
	if err := models.DB.Create(newSkill).Error; err != nil {
		return nil, err
	}
	return newSkill, nil
}

func (repo *SkillRepository) Count() (int, error) {
	var count int64
	query := models.DB.Model(&models.Skill{})
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *SkillRepository) ReadAll() ([]models.Skill, error) {
	var skills []models.Skill
	if err := models.DB.Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func (repo *SkillRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.Skill, error) {
	var skills []models.Skill

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
	if err := query.Offset(offset).Limit(pageSize).Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func (repo *SkillRepository) Read(id int64) (*models.Skill, error) {
	var skill models.Skill
	if err := models.DB.First(&skill, id).Error; err != nil {
		return nil, err
	}
	return &skill, nil
}

func (repo *SkillRepository) ReadByCode(code string) (*models.Skill, error) {
	var skill models.Skill
	if err := models.DB.Where("code = ?", code).First(&skill).Error; err != nil {
		return nil, err
	}
	return &skill, nil
}

func (repo *SkillRepository) ReadByName(name string) (*models.Skill, error) {
	var skill models.Skill
	if err := models.DB.Where("name = ?", name).First(&skill).Error; err != nil {
		return nil, err
	}
	return &skill, nil
}

func (repo *SkillRepository) ReadByNameOrCode(name string, code string) (*models.Skill, error) {
	var skill models.Skill
	if err := models.DB.Where("name = ? OR code = ?", name, code).First(&skill).Error; err != nil {
		return nil, err
	}
	return &skill, nil
}
