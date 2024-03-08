package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserAttachmentRepository struct{}

func NewUserAttachmentRepository() *UserAttachmentRepository {
	return &UserAttachmentRepository{}
}

func (repo *UserAttachmentRepository) Create(newData *models.UserAttachment) (*models.UserAttachment, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
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
	var datas []models.UserAttachment

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if category != nil {
		query = query.Where("category LIKE ?", category)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserAttachmentRepository) ReadFilteredPaginated(userId *int64, category *string, pageSize, pageNumber int) ([]models.UserAttachment, error) {
	var datas []models.UserAttachment

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
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserAttachmentRepository) Read(id int64) (*models.UserAttachment, error) {
	var data models.UserAttachment
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserAttachmentRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserAttachment{}, id).Error; err != nil {
		return err
	}
	return nil
}
