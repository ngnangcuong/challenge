package usecase

import (
	"testing"

	"challenge3/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
)

type MockUserRepo struct {
	mock.Mock
}

func (mock *MockUserRepo) Select() ([]models.User, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]models.User), args.Error(1)
}

func (mock *MockUserRepo) Find(email string) (models.User, error) {
	args := mock.Called(email)
	result := args.Get(0)
	return result.(models.User), args.Error(1)
}

func (mock *MockUserRepo) Create(user models.User) (error) {
	args := mock.Called(user)
	return args.Error(0)
}

func (mock *MockUserRepo) Insert(user models.User) (error) {
	args := mock.Called(user)
	return args.Error(0)
}

func (mock *MockUserRepo) Update(user models.User) (error) {
	args := mock.Called(user)
	return args.Error(0)
}

func (mock *MockUserRepo) Delete(email string) (error) {
	args := mock.Called(email)
	return args.Error(0)
}

func TestFindUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	var (
		user = models.User{
			Email: "nangcuong",
			Name: "Cuong",
			Password: "pa$$w0rd",
			Role: "user",
		}
		email = "nangcuong"
		emailFake = "nanguong"
	)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find", email).Return(user, nil).Once()

		userService := NewUserService(mockRepo)
		result, err := userService.FindUser(email)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Does-not-exist-user", func(t *testing.T) {
		mockRepo.On("Find", emailFake).Return(models.User{}, gorm.ErrRecordNotFound)

		userService := NewUserService(mockRepo)
		result, err := userService.FindUser(emailFake)

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, models.User{}, result)
	})
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	var (
		email = "nangcuong"
		name = "Cuong"
		password = "pa$$w0rd"
		user = models.User{
			Email: "nangcuong",
			Name: "Cuong",
			Password: "pa$$w0rd",
			Role: "user",
		}
	)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find", email).Return(models.User{}, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", user).Return(nil).Once()

		userService := NewUserService(mockRepo)
		result, err := userService.CreateUser(email, name, password)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Users-already-exist", func(t *testing.T) {
		mockRepo.On("Find", email).Return(user, nil).Once()

		userService := NewUserService(mockRepo)
		result, err := userService.CreateUser(email, name, password)

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, models.User{}, result)
	})
}

func TestGetListUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	var listUser = []models.User{
		0: {
			Email: "admin",
			Name: "Admin",
			Password: "admin",
			Role: "admin",
		},
		1: {
			Email: "nangcuong",
			Name: "Cuong",
			Password: "pa$$w0rd",
			Role: "user",
		},
		2: {
			Email: "thamnh",
			Name: "Tham",
			Password: "040701",
			Role: "Police",
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Select").Return(listUser, nil).Once()

		userService := NewUserService(mockRepo)
		userList, err := userService.GetListUser()

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, listUser, userList)
	})

	t.Run("Failed", func(t *testing.T) {
		mockRepo.On("Select").Return([]models.User(nil), gorm.ErrRecordNotFound).Once()

		userService := NewUserService(mockRepo)
		userList, err := userService.GetListUser()

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, []models.User(nil), userList)
	})
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	var (
		user = models.User{
			Email: "nangcuong",
			Name: "Cuong",
			Password: "pa$$w0rd",
			Role: "user",
		}
		email = "nangcuong"
		emailFake = "nanguong"
	)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find", email).Return(user, nil).Once()
		mockRepo.On("Delete", email).Return(nil).Once()

		userService := NewUserService(mockRepo)
		err := userService.DeleteUser(email)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Does-not-exist-user", func(t *testing.T) {
		mockRepo.On("Find", emailFake).Return(models.User{}, gorm.ErrRecordNotFound)
		
		userService := NewUserService(mockRepo)
		err := userService.DeleteUser(emailFake)

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	var (
		user = models.User{
			Email: "nangcuong",
			Name: "Cuong",
			Password: "pa$$w0rd",
			Role: "user",
		}
		userFake = models.User{
			Email: "nanguong",
			Name: "NangCuong",
			Password: "password",
			Role: "Police",
		}
		userUpdate = models.User{
			Email: "nangcuong",
			Name: "NangCuong",
			Password: "password",
			Role: "user",
		}
		email = "nangcuong"
		emailFake = "nanguong"
	)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find", email).Return(user, nil).Once()
		mockRepo.On("Update", userUpdate).Return(nil).Once()

		userService := NewUserService(mockRepo)
		err := userService.UpdateUser(userUpdate)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Does-not-exist-user", func(t *testing.T) {
		mockRepo.On("Find", emailFake).Return(models.User{}, gorm.ErrRecordNotFound).Once()

		userService := NewUserService(mockRepo)
		err := userService.UpdateUser(userFake)

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
	})
}	
