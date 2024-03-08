package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type SkillRepository struct{}

func NewSkillRepository() *SkillRepository {
	return &SkillRepository{}
}

func (repo *SkillRepository) Create(newData *models.Skill) (*models.Skill, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
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
	var datas []models.Skill
	if err := models.DB.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *SkillRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.Skill, error) {
	var datas []models.Skill

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

func (repo *SkillRepository) Read(id int64) (*models.Skill, error) {
	var data models.Skill
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *SkillRepository) ReadByCode(code string) (*models.Skill, error) {
	var data models.Skill
	if err := models.DB.Where("code = ?", code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *SkillRepository) ReadByName(name string) (*models.Skill, error) {
	var data models.Skill
	if err := models.DB.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *SkillRepository) ReadByNameOrCode(name string, code string) (*models.Skill, error) {
	var data models.Skill
	if err := models.DB.Where("name = ? OR code = ?", name, code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
