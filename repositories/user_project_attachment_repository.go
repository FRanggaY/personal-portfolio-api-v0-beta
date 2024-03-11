package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserProjectAttachmentRepository struct{}

func NewUserProjectAttachmentRepository() *UserProjectAttachmentRepository {
	return &UserProjectAttachmentRepository{}
}

func (repo *UserProjectAttachmentRepository) Create(newData *models.UserProjectAttachment) (*models.UserProjectAttachment, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserProjectAttachmentRepository) Count() (int, error) {
	var count int64
	query := models.DB.Model(&models.UserProjectAttachment{})

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserProjectAttachmentRepository) ReadAll(userProjectID *int64) ([]models.UserProjectAttachment, error) {
	query := models.DB
	var datas []models.UserProjectAttachment

	if userProjectID != nil {
		query = query.Where("user_project_id", userProjectID)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserProjectAttachmentRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.UserProjectAttachment, error) {
	var datas []models.UserProjectAttachment

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

func (repo *UserProjectAttachmentRepository) Read(ID int64) (*models.UserProjectAttachment, error) {
	var data models.UserProjectAttachment
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserProjectAttachmentRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserProjectAttachment{}, ID).Error; err != nil {
		return err
	}
	return nil
}
