package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type UserProjectRepository struct{}

func NewUserProjectRepository() *UserProjectRepository {
	return &UserProjectRepository{}
}

func (repo *UserProjectRepository) Create(newData *models.UserProject) (*models.UserProject, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *UserProjectRepository) Count(userID *int64, projectPlatformID *int64, isActive *bool) (int, error) {
	var count int64
	query := models.DB.Model(&models.UserProject{})

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if projectPlatformID != nil {
		query = query.Where("project_platform_id", projectPlatformID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserProjectRepository) ReadAll(userID *int64, projectPlatformID *int64, isActive *bool) ([]models.UserProject, error) {
	query := models.DB
	var datas []models.UserProject

	if userID != nil {
		query = query.Where(helper.FilterUserIDEqual, userID)
	}

	if projectPlatformID != nil {
		query = query.Where("project_platform_id", projectPlatformID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserProjectRepository) ReadFilteredPaginated(userID *int64, projectPlatformID *int64, isActive *bool, pageSize, pageNumber int) ([]models.UserProject, error) {
	var datas []models.UserProject

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

	if projectPlatformID != nil {
		query = query.Where("project_platform_id", projectPlatformID)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", isActive)
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserProjectRepository) Read(ID int64) (*models.UserProject, error) {
	var data models.UserProject
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserProjectRepository) ReadByUserIDSlugLanguageID(userID int64, slug string, languageID int64) (*models.ProjectTranslationResponse, error) {
	var data models.ProjectTranslationResponse
	if err := models.DB.
		Table("user_project_translations").
		Select(`
		user_project_translations.*,
		user_projects.image_url,
		user_projects.slug,
		user_projects.project_created_at,
		user_projects.project_updated_at,
		user_projects.project_platform_id
	`).
		Joins("LEFT JOIN user_projects ON user_project_translations.user_project_id = user_projects.id").
		Where("user_projects.user_id = ? AND user_projects.slug = ?", userID, slug).
		Where("user_project_translations.language_id = ?", languageID).
		First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserProjectRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.UserProject{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserProjectRepository) ReadTranslationsByUserIDLanguageID(userID int64, projectPlatformID *int64, languageID int64, isActive *bool, pageNumber int, pageSize int) ([]models.ProjectTranslationResponse, error) {
	var skills []models.ProjectTranslationResponse

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
		Table("user_project_translations").
		Select(`
			user_project_translations.*,
			user_projects.image_url,
			user_projects.slug,
			user_projects.project_created_at,
			user_projects.project_updated_at,
			user_projects.project_platform_id
        `).
		Joins("LEFT JOIN user_projects ON user_project_translations.user_project_id = user_projects.id").
		Where("user_projects.user_id = ?", userID).
		Where("user_project_translations.language_id = ?", languageID).
		Limit(pageSize).Offset(offset)

	if isActive != nil {
		query = query.Where("user_projects.is_active = ?", isActive)
	}

	if projectPlatformID != nil {
		query = query.Where("user_projects.project_platform_id = ?", projectPlatformID)
	}

	if err := query.Find(&skills).Error; err != nil {
		return nil, err
	}

	return skills, nil
}
