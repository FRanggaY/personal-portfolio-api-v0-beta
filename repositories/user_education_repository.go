package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserEducationRepository struct{}

func NewUserEducationRepository() *UserEducationRepository {
	return &UserEducationRepository{}
}

func (repo *UserEducationRepository) Create(newData *models.UserEducation) (*models.UserEducation, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserEducationRepository) Count(userId *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserEducation{})

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserEducationRepository) ReadAll(userId *int64) ([]models.UserEducation, error) {
	query := models.DB
	var datas []models.UserEducation

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserEducationRepository) ReadFilteredPaginated(userId *int64, pageSize, pageNumber int) ([]models.UserEducation, error) {
	var datas []models.UserEducation

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

func (repo *UserEducationRepository) Read(id int64) (*models.UserEducation, error) {
	var data models.UserEducation
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserEducationRepository) ReadByUserIdSchoolId(userId int64, schoolId int64) (*models.UserEducation, error) {
	var data models.UserEducation
	if err := models.DB.Where("user_id = ? AND school_id = ?", userId, schoolId).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserEducationRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserEducation{}, id).Error; err != nil {
		return err
	}
	return nil
}
