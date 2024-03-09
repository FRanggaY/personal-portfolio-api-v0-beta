package repositories

import (
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

func (repo *UserExperienceRepository) Count(userId *int64) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserExperience{})

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserExperienceRepository) ReadAll(userId *int64) ([]models.UserExperience, error) {
	query := models.DB
	var datas []models.UserExperience

	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserExperienceRepository) ReadFilteredPaginated(userId *int64, pageSize, pageNumber int) ([]models.UserExperience, error) {
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
	if userId != nil {
		query = query.Where("user_id LIKE ?", userId)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserExperienceRepository) Read(id int64) (*models.UserExperience, error) {
	var data models.UserExperience
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserExperienceRepository) ReadByUserIdCompanyId(userId int64, companyId int64) (*models.UserExperience, error) {
	var data models.UserExperience
	if err := models.DB.Where("user_id = ? AND company_id = ?", userId, companyId).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserExperienceRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.UserExperience{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserExperienceRepository) ReadTranslationsByUserIDLanguageID(userID int64, languageID int64, pageNumber int, pageSize int) ([]models.ExperienceTranslationResponse, error) {
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

	if err := query.Find(&skills).Error; err != nil {
		return nil, err
	}

	return skills, nil
}
