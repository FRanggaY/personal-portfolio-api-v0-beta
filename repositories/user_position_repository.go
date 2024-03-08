package repositories

import (
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

func (repo *UserPositionRepository) Count(userId *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserPosition{})

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserPositionRepository) ReadAll(userId *int64) ([]models.UserPosition, error) {
	query := models.DB
	var datas []models.UserPosition

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserPositionRepository) ReadFilteredPaginated(userId *int64, pageSize, pageNumber int) ([]models.UserPosition, error) {
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
	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserPositionRepository) Read(id int64) (*models.UserPosition, error) {
	var data models.UserPosition
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserPositionRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserPosition{}, id).Error; err != nil {
		return err
	}
	return nil
}
