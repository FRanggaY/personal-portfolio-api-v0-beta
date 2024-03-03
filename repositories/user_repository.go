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

func (repo *UserRepository) Count(nameFilter string) (int, error) {
	var count int64
	query := models.DB.Model(&models.User{})
	if nameFilter != "" {
		query = query.Where("name LIKE ?", "%"+nameFilter+"%")
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

	newUser := models.User{
		Username: userInput.Username,
		Name:     userInput.Name,
		Password: hashedPassword,
	}

	// Insert user into database
	if err := models.DB.Create(&newUser).Error; err != nil {
		return nil, err
	}
	return &newUser, nil
}

func (repo *UserRepository) ReadAll() ([]models.User, error) {
	var users []models.User
	if err := models.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepository) ReadFilteredPaginated(nameFilter string, pageSize, pageNumber int) ([]models.User, error) {
	var users []models.User

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
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepository) Read(id int64) (*models.User, error) {
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) ReadByUsername(username string) (*models.User, error) {
	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) Update(id int64, updatedUser *models.UserEditForm) error {
	existingUser, err := repo.Read(id)
	if err != nil {
		return err
	}

	existingUser.Name = updatedUser.Name
	existingUser.Username = updatedUser.Username

	if err := models.DB.Save(existingUser).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(id int64) error {
	if err := models.DB.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
