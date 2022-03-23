package usecase

import (
	"testing"

	"challenge3/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
)

type MockRoleRepo struct {
	mock.Mock
}

func (mock *MockRoleRepo) Create(role models.Role) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *MockRoleRepo) Find(name string) (models.Role, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(models.Role), args.Error(1)

}

func TestFind(t *testing.T) {
	mockRepo := new(MockRoleRepo)

	var roleCheck = models.Role{
		Name: "Police",
		Permission: "rd",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find").Return(roleCheck, nil).Once()
		
		roleService := NewRoleService(mockRepo)

		result, err := roleService.Find("Police")

		mockRepo.AssertExpectations(t)
		assert.Equal(t, roleCheck, result)
		assert.NoError(t, err)
	})

	t.Run("Does-not-exist-role", func(t *testing.T) {
		mockRepo.On("Find").Return(models.Role{}, gorm.ErrRecordNotFound).Once()

		roleService := NewRoleService(mockRepo)
		result, err := roleService.Find("Support")

		mockRepo.AssertExpectations(t)
		assert.Equal(t, models.Role{}, result)
		assert.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockRoleRepo)

	var role = models.Role{
		Name: "Police",
		Permission: "rd",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find").Return(models.Role{}, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create").Return(nil).Once()

		roleService := NewRoleService(mockRepo)
		err := roleService.Create("Police", "rd")

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Role-is-already-existed", func(t *testing.T) {
		mockRepo.On("Find").Return(role, nil).Once()
		//mockRepo.On("Create").Return(nil).Once()

		roleService := NewRoleService(mockRepo)
		err := roleService.Create("Police", "rd")

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
	})
}


