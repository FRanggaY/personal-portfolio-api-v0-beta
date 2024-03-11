package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserProjectRepository struct{}

func NewUserProjectRepository() *UserProjectRepository {
	return &UserProjectRepository{}
}

func (repo *UserProjectRepository) Create(newData *models.UserProject) (*models.UserProject, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserProjectRepository) Count(userID *int64, projectPlatformID *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserProject{})

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if projectPlatformID != nil {
		query = query.Where("project_platform_id", projectPlatformID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserProjectRepository) ReadAll(userID *int64, projectPlatformID *int64, isActive *bool) ([]models.UserProject, error) {
	query := models.DB
	var datas []models.UserProject

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if projectPlatformID != nil {
		query = query.Where("project_platform_id", projectPlatformID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserProjectRepository) ReadFilteredPaginated(userID *int64, projectPlatformID *int64, pageSize, pageNumber int) ([]models.UserProject, error) {
	var datas []models.UserProject

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

	if projectPlatformID != nil {
		query = query.Where("project_platform_id", projectPlatformID)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserProjectRepository) Read(ID int64) (*models.UserProject, error) {
	var data models.UserProject
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserProjectRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserProject{}, ID).Error; err != nil {
		return err
	}
	return nil
}
