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

func (repo *UserEducationRepository) ReadTranslationsByUserIDLanguageID(userID int64, languageID int64, pageNumber int, pageSize int) ([]models.EducationTranslationResponse, error) {
	var skills []models.EducationTranslationResponse

	// default
	if pageSize <= 0 {
		pageSize = 5
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	// calculate offset
	offset := (pageNumber - 1) * pageSize

	query := models.DB.
		Table("user_education_translations").
		Select(`
			user_education_translations.*,
			IFNULL(schools.id, '') AS school_id, 
            IFNULL(schools.code, '') AS school_code,
            IFNULL(schools.name, '') AS school_name,
            IFNULL(schools.image_url, '') AS school_image_url,
            IFNULL(schools.url, '') AS school_url,
            IFNULL(schools.is_external_url, '') AS school_is_external_url,
            IFNULL(schools.is_external_image_url, '') AS school_is_external_image_url,
            IFNULL(schools.address, '') AS school_address
        `).
		Joins("LEFT JOIN user_educations ON user_education_translations.user_education_id = user_educations.id").
		Joins("LEFT JOIN schools ON user_educations.school_id = schools.id").
		Where("user_educations.user_id = ?", userID).
		Where("user_education_translations.language_id = ?", languageID).
		Limit(pageSize).Offset(offset)

	if err := query.Find(&skills).Error; err != nil {
		return nil, err
	}

	return skills, nil
}
