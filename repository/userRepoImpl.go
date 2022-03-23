package repository

import (
	"challenge3/models"
	"github.com/jinzhu/gorm"
)

type userRepoImpl struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) models.UserRepo {
	return &userRepoImpl{
		DB: db,
	}
}

func (u *userRepoImpl) Select() ([]models.User, error) {
	var userList []models.User

	result := u.DB.Find(&userList)

	if result.Error != nil {
		return nil, result.Error
	}
	return userList, nil
}

func (u *userRepoImpl) Insert(user models.User) (error) {
	result := u.DB.Create(&user)

	return result.Error
}

func (u *userRepoImpl) Update(user models.User) (error) {

	result := u.DB.Save(&user)

	return result.Error
}

func (u *userRepoImpl) Delete(email string) (error) {
	user, err := u.Find(email)
	if err != nil {
		return err
	}

	result := u.DB.Delete(&user)
	return result.Error
}

func (u *userRepoImpl) Find(email string) (models.User, error) {
	var user models.User

	result := u.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func (u *userRepoImpl) Create(user models.User) error {
	result := u.DB.Create(&user)
	return result.Error
}