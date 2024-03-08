package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserExperienceRepository struct{}

func NewUserExperienceRepository() *UserExperienceRepository {
	return &UserExperienceRepository{}
}

func (repo *UserExperienceRepository) Create(newData *models.UserExperience) (*models.UserExperience, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserExperienceRepository) Count(userId *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserExperience{})

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserExperienceRepository) ReadAll(userId *int64) ([]models.UserExperience, error) {
	query := models.DB
	var datas []models.UserExperience

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserExperienceRepository) ReadFilteredPaginated(userId *int64, pageSize, pageNumber int) ([]models.UserExperience, error) {
	var datas []models.UserExperience

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
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserExperienceRepository) Read(id int64) (*models.UserExperience, error) {
	var data models.UserExperience
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserExperienceRepository) ReadByUserIdCompanyId(userId int64, companyId int64) (*models.UserExperience, error) {
	var data models.UserExperience
	if err := models.DB.Where("user_id = ? AND company_id = ?", userId, companyId).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserExperienceRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserExperience{}, id).Error; err != nil {
		return err
	}
	return nil
}
