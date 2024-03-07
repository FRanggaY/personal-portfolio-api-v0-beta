package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserPositionRepository struct{}

func NewUserPositionRepository() *UserPositionRepository {
	return &UserPositionRepository{}
}

func (repo *UserPositionRepository) Create(newUserPosition *models.UserPosition) (*models.UserPosition, error) {
	// Insert skill into database
	if err := models.DB.Create(newUserPosition).Error; err != nil {
		return nil, err
	}
	return newUserPosition, nil
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
	var userPositions []models.UserPosition

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Find(&userPositions).Error; err != nil {
		return nil, err
	}
	return userPositions, nil
}

func (repo *UserPositionRepository) ReadFilteredPaginated(userId *int64, pageSize, pageNumber int) ([]models.UserPosition, error) {
	var userPositions []models.UserPosition

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
	if err := query.Offset(offset).Limit(pageSize).Find(&userPositions).Error; err != nil {
		return nil, err
	}
	return userPositions, nil
}

func (repo *UserPositionRepository) Read(id int64) (*models.UserPosition, error) {
	var userPosition models.UserPosition
	if err := models.DB.First(&userPosition, id).Error; err != nil {
		return nil, err
	}
	return &userPosition, nil
}

func (repo *UserPositionRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserPosition{}, id).Error; err != nil {
		return err
	}
	return nil
}
