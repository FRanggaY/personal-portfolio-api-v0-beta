package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserPositionRepository struct{}

func NewUserPositionRepository() *UserPositionRepository {
	return &UserPositionRepository{}
}

func (repo *UserPositionRepository) Create(newData *models.UserPosition) (*models.UserPosition, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserPositionRepository) Count(userID *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserPosition{})

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserPositionRepository) ReadAll(userID *int64, isActive *bool) ([]models.UserPosition, error) {
	query := models.DB
	var datas []models.UserPosition

	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserPositionRepository) ReadFilteredPaginated(userID *int64, pageSize, pageNumber int) ([]models.UserPosition, error) {
	var datas []models.UserPosition

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
		query = query.Where("user_id LIKE ?", userID)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserPositionRepository) Read(ID int64) (*models.UserPosition, error) {
	var data models.UserPosition
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserPositionRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserPosition{}, ID).Error; err != nil {
		return err
	}
	return nil
}
