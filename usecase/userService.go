package usecase

import (
	"challenge3/models"
)

type UserService struct {
	repo models.UserRepo
}

func NewUserUseCase(repo models.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) FindUser(email string) (models.User, error) {
	return u.repo.Find(email)
}

func (u *UserService) CreateUser(email, name, password string) (models.User, error) {
	var user = models.User{
		Email: email,
		Name: name,
		Password: password,
		Role: "user",
	}

	err := u.repo.Create(user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UserService) GetListUser() ([]models.User, error) {
	return u.repo.Select()
}

func (u *UserService) DeleteUser(email string) (error) {
	_, err := u.repo.Find(email)
	if err != nil {
		return err
	}

	return u.repo.Delete(email)
}

func (u *UserService) UpdateUser(user models.User) (error) {
	user, err := u.repo.Find(user.Email)
	if err != nil {
		return err
	}

	return u.repo.Update(user)
}
