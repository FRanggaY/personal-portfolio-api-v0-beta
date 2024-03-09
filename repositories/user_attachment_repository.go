package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
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

func (repo *UserAttachmentRepository) Count(userID *int64, category *string) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserAttachment{})

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if category != nil {
		query = query.Where(helper.FilterCategoryLike, category)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserAttachmentRepository) ReadAll(userID *int64, category *string) ([]models.UserAttachment, error) {
	query := models.DB
	var datas []models.UserAttachment

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if category != nil {
		query = query.Where(helper.FilterCategoryLike, category)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserAttachmentRepository) ReadFilteredPaginated(userID *int64, category *string, pageSize, pageNumber int) ([]models.UserAttachment, error) {
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
	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if category != nil {
		query = query.Where(helper.FilterCategoryLike, category)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserAttachmentRepository) Read(ID int64) (*models.UserAttachment, error) {
	var data models.UserAttachment
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserAttachmentRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserAttachment{}, ID).Error; err != nil {
		return err
	}
	return nil
}
