package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserExperienceRepository struct{}

func NewUserExperienceRepository() *UserExperienceRepository {
	return &UserExperienceRepository{}
}

func (repo *UserExperienceRepository) Create(newData *models.UserExperience) (*models.UserExperience, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserExperienceRepository) Count(userID *int64, isActive *bool) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserExperience{})

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserExperienceRepository) ReadAll(userID *int64) ([]models.UserExperience, error) {
	query := models.DB
	var datas []models.UserExperience

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserExperienceRepository) ReadFilteredPaginated(userID *int64, pageSize, pageNumber int) ([]models.UserExperience, error) {
	var datas []models.UserExperience

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

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserExperienceRepository) Read(ID int64) (*models.UserExperience, error) {
	var data models.UserExperience
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserExperienceRepository) ReadByUserIDCompanyID(userID int64, companyID int64) (*models.UserExperience, error) {
	var data models.UserExperience
	if err := models.DB.Where("user_id = ? AND company_id = ?", userID, companyID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserExperienceRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserExperience{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserExperienceRepository) DeleteByUserIDCompanyID(userID int64, companyID int64) error {
	if err := models.DB.
		Where("user_id = ? AND company_id = ?", userID, companyID).
		Delete(&models.UserExperience{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserExperienceRepository) ReadTranslationsByUserIDLanguageID(userID int64, languageID int64, isActive *bool, pageNumber int, pageSize int) ([]models.ExperienceTranslationResponse, error) {
	var skills []models.ExperienceTranslationResponse

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
		Table("user_experience_translations").
		Select(`
			user_experience_translations.*,
			user_experiences.month_start,
			user_experiences.month_end,
			user_experiences.year_start,
			user_experiences.year_end,
			IFNULL(companies.id, '') AS company_id, 
            IFNULL(companies.code, '') AS company_code,
            IFNULL(companies.name, '') AS company_name,
            IFNULL(companies.image_url, '') AS company_image_url,
            IFNULL(companies.url, '') AS company_url,
            IFNULL(companies.is_external_url, '') AS company_is_external_url,
            IFNULL(companies.is_external_image_url, '') AS company_is_external_image_url,
            IFNULL(companies.address, '') AS company_address
        `).
		Joins("LEFT JOIN user_experiences ON user_experience_translations.user_experience_id = user_experiences.id").
		Joins("LEFT JOIN companies ON user_experiences.company_id = companies.id").
		Where("user_experiences.user_id = ?", userID).
		Where("user_experience_translations.language_id = ?", languageID).
		Limit(pageSize).Offset(offset)

	if isActive != nil {
		query = query.Where("user_experiences.is_active = ?", isActive)
	}

	if err := query.Find(&skills).Error; err != nil {
		return nil, err
	}

	return skills, nil
}
