package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type CompanyRepository struct{}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{}
}

func (repo *CompanyRepository) Create(newData *models.Company) (*models.Company, error) {
	if err := models.DB.Create(newData).Error; err != nil {
		return nil, err
	}
	return newData, nil
}

func (repo *CompanyRepository) Count() (int, error) {
	var count int64
	query := models.DB.Model(&models.Company{})
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *CompanyRepository) ReadAll() ([]models.Company, error) {
	var datas []models.Company
	if err := models.DB.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *CompanyRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.Company, error) {
	var datas []models.Company

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

func (repo *CompanyRepository) Read(id int64) (*models.Company, error) {
	var data models.Company
	if err := models.DB.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *CompanyRepository) ReadByCode(code string) (*models.Company, error) {
	var data models.Company
	if err := models.DB.Where("code = ?", code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *CompanyRepository) ReadByName(name string) (*models.Company, error) {
	var data models.Company
	if err := models.DB.Where("name = ?", name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *CompanyRepository) ReadByNameOrCode(name string, code string) (*models.Company, error) {
	var data models.Company
	if err := models.DB.Where("name = ? OR code = ?", name, code).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
