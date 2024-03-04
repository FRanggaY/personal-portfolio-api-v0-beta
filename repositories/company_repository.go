package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

type CompanyRepository struct{}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{}
}

func (repo *CompanyRepository) Create(newCompany *models.Company) (*models.Company, error) {
	// Insert Company into database
	if err := models.DB.Create(newCompany).Error; err != nil {
		return nil, err
	}
	return newCompany, nil
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
	var Companys []models.Company
	if err := models.DB.Find(&Companys).Error; err != nil {
		return nil, err
	}
	return Companys, nil
}

func (repo *CompanyRepository) ReadFilteredPaginated(pageSize, pageNumber int) ([]models.Company, error) {
	var Companys []models.Company

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
	if err := query.Offset(offset).Limit(pageSize).Find(&Companys).Error; err != nil {
		return nil, err
	}
	return Companys, nil
}

func (repo *CompanyRepository) Read(id int64) (*models.Company, error) {
	var Company models.Company
	if err := models.DB.First(&Company, id).Error; err != nil {
		return nil, err
	}
	return &Company, nil
}
