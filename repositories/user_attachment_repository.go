package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserAttachmentRepository struct{}

func NewUserAttachmentRepository() *UserAttachmentRepository {
	return &UserAttachmentRepository{}
}

func (repo *UserAttachmentRepository) Create(newUserAttachment *models.UserAttachment) (*models.UserAttachment, error) {
	// Insert skill into database
	if err := models.DB.Create(newUserAttachment).Error; err != nil {
		return nil, err
	}
	return newUserAttachment, nil
}

func (repo *UserAttachmentRepository) Count(userId *int64, category *string) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserAttachment{})

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if category != nil {
		query = query.Where("category LIKE ?", category)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserAttachmentRepository) ReadAll(userId *int64, category *string) ([]models.UserAttachment, error) {
	query := models.DB
	var userAttachments []models.UserAttachment

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if category != nil {
		query = query.Where("category LIKE ?", category)
	}

	if err := query.Find(&userAttachments).Error; err != nil {
		return nil, err
	}
	return userAttachments, nil
}

func (repo *UserAttachmentRepository) ReadFilteredPaginated(userId *int64, category *string, pageSize, pageNumber int) ([]models.UserAttachment, error) {
	var userAttachments []models.UserAttachment

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

	if category != nil {
		query = query.Where("category LIKE ?", category)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&userAttachments).Error; err != nil {
		return nil, err
	}
	return userAttachments, nil
}

func (repo *UserAttachmentRepository) Read(id int64) (*models.UserAttachment, error) {
	var userAttachment models.UserAttachment
	if err := models.DB.First(&userAttachment, id).Error; err != nil {
		return nil, err
	}
	return &userAttachment, nil
}

func (repo *UserAttachmentRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserAttachment{}, id).Error; err != nil {
		return err
	}
	return nil
}
