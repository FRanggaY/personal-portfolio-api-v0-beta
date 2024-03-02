package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
)

func CreateSchool(newSchool *models.School) error {
	// Insert school into database
	if err := models.DB.Create(&newSchool).Error; err != nil {
		return err
	}
	return nil
}

func CountSchool() (int, error) {
	var count int64
	query := models.DB.Model(&models.School{})
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func ReadAllSchool() ([]models.School, error) {
	var schools []models.School
	if err := models.DB.Find(&schools).Error; err != nil {
		return nil, err
	}
	return schools, nil
}

func ReadSchoolsFilteredPaginated(pageSize, pageNumber int) ([]models.School, error) {
	var schools []models.School

	// default
	if pageSize <= 0 {
		pageSize = 5
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	// calculate off set
	offset := (pageNumber - 1) * pageSize

	query := models.DB

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&schools).Error; err != nil {
		return nil, err
	}
	return schools, nil
}

func ReadSchool(id int64) (*models.School, error) {
	var school models.School
	if err := models.DB.First(&school, id).Error; err != nil {
		return nil, err
	}
	return &school, nil
}
