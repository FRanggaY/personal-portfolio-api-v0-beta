package repositories

import (
	"github.com/FRanggaY/personal-portfolio-api/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) HashUserPassword(password string) (string, error) {
	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (repo *UserRepository) CompareUserPassword(hashedPassword, plainPassword string) error {
	// Compare hashed password with plain password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) Count(name string) (int, error) {
	var count int64
	query := models.DB.Model(&models.User{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserRepository) Create(userInput *models.UserCreateForm) (*models.User, error) {
	// Hash user's password
	hashedPassword, err := repo.HashUserPassword(userInput.Password)
	if err != nil {
		return nil, err
	}

	newData := models.User{
		Username: userInput.Username,
		Name:     userInput.Name,
		Password: hashedPassword,
	}

	// Insert user into database
	if err := models.DB.Create(&newData).Error; err != nil {
		return nil, err
	}
	return &newData, nil
}

func (repo *UserRepository) ReadAll() ([]models.User, error) {
	var datas []models.User
	if err := models.DB.Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserRepository) ReadFilteredPaginated(nameFilter string, pageSize, pageNumber int) ([]models.User, error) {
	var datas []models.User

	// default
	if pageSize <= 0 {
		pageSize = 5
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	// calculate off set
	offset := (pageNumber - 1) * pageSize

	// filter name
	query := models.DB
	if nameFilter != "" {
		query = query.Where("name LIKE ?", "%"+nameFilter+"%")
	}

	// pagination
	if err := query.Offset(offset).Limit(pageSize).Find(&datas).Error; err != nil {
		return nil, err
	}
	return datas, nil
}

func (repo *UserRepository) Read(ID int64) (*models.User, error) {
	var data models.User
	if err := models.DB.First(&data, ID).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserRepository) ReadByUsername(username string) (*models.User, error) {
	var data models.User
	if err := models.DB.Where("username = ?", username).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *UserRepository) Update(ID int64, updatedUser *models.UserEditForm) error {
	existingData, err := repo.Read(ID)
	if err != nil {
		return err
	}

	existingData.Name = updatedUser.Name
	existingData.Username = updatedUser.Username

	if err := models.DB.Save(existingData).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(ID int64) error {
	if err := models.DB.Delete(&models.User{}, ID).Error; err != nil {
		return err
	}
	return nil
}
