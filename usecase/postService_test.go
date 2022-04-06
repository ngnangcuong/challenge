package usecase

import (
	"testing"

	"challenge3/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
)

type MockPostRepo struct {
	mock.Mock
}

type MockPostSearchRepo struct {
	mock.Mock
}

func (mock *MockPostSearchRepo) Search(keyword string) ([]models.Post, error) {
	args := mock.Called(keyword)
	result := args.Get(0)
	return result.([]models.Post), args.Error(1)
}

func (mock *MockPostSearchRepo)	Index(post models.Post) error {
	args := mock.Called(post)
	return args.Error(0)
}

func (mock *MockPostSearchRepo)	Update(post models.Post) error {
	args := mock.Called(post)
	return args.Error(0)
}

func (mock *MockPostSearchRepo)	Delete(id string) error {
	args := mock.Called(id)
	return args.Error(0)
}

func (mock *MockPostRepo) Select() ([]models.Post, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]models.Post), args.Error(1)
}

func (mock *MockPostRepo) Delete(id uint) (error) {
	args := mock.Called(id)
	return args.Error(0)
}

func (mock *MockPostRepo) Update(post models.Post)	(error) {
	args := mock.Called(post)
	return args.Error(0)
}

func (mock *MockPostRepo) Create(post models.Post) (models.Post, error) {
	args := mock.Called(post)
	result := args.Get(0)
	return result.(models.Post), args.Error(1)
}

func (mock *MockPostRepo) Find(id uint) (models.Post, error) {
	args := mock.Called(id)
	result := args.Get(0)
	return result.(models.Post), args.Error(1)
}

func TestGetListPost(t *testing.T) {
	mockRepo := new(MockPostRepo)
	mockSearchRepo := new(MockPostSearchRepo)
	var listPost = []models.Post{
		0: {
			ID: 1,
			UserID: 1,
			Email: "admin",
			Content: "This is admin's post",
		},
		1: {
			ID: 2,
			UserID: 2,
			Email: "nangcuong",
			Content: "This is nangcuong's post",
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Select").Return(listPost, nil).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		postList, err := postService.GetListPost()

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, listPost, postList)
	})

	t.Run("Empty-postList", func(t *testing.T) {
		mockRepo.On("Select").Return([]models.Post(nil), gorm.ErrRecordNotFound).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		postList, err := postService.GetListPost()

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, postList)
	})
}

func TestCreatePost(t *testing.T) {
	mockSearchRepo := new(MockPostSearchRepo)
	mockRepo := new(MockPostRepo)
	var post = models.Post{
		ID: 1,
		UserID: 1,
		Email: "admin",
		Content: "This is admin's post",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Create", post).Return(post, nil).Once()
		mockSearchRepo.On("Index", post).Return(nil).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		result, err := postService.CreatePost(post)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, post, result)
	})
}

func TestFindPost(t *testing.T) {
	mockRepo := new(MockPostRepo)
	mockSearchRepo := new(MockPostSearchRepo)
	var (
		post = models.Post{
			ID: 10,
			UserID: 1,
			Email: "admin",
			Content: "Another post from admin",
		}
		id = uint(10)
		idFake = uint(7)
	) 


	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find", id).Return(post, nil).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		result, err := postService.Find(id)

		mockRepo.AssertExpectations(t)
		assert.NoError(t ,err)
		assert.Equal(t, post, result)
	})

	t.Run("Does-not-exist-post", func(t *testing.T) {
		mockRepo.On("Find", idFake).Return(models.Post{}, gorm.ErrRecordNotFound).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		result, err := postService.Find(idFake)

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, models.Post{}, result)
	})
}

func TestDeletePost(t *testing.T) {
	mockRepo := new(MockPostRepo)
	mockSearchRepo := new(MockPostSearchRepo)
	var (
		post = models.Post{
			ID: 7,
			UserID: 7,
			Email: "nangcuong",
			Content: "Posted by nangcuong",
		}
		id = uint(7)
		idES = "7"
		idFake = uint(14)
	)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find", id).Return(post, nil).Once()
		mockRepo.On("Delete", id).Return(nil).Once()
		mockSearchRepo.On("Delete", idES).Return(nil).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		err := postService.DeletePost(id)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Does-not-exist-post", func(t *testing.T) {
		mockRepo.On("Find", idFake).Return(models.Post{}, gorm.ErrRecordNotFound).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		err := postService.DeletePost(idFake)

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
	})
}

func TestUpdatePost(t *testing.T) {
	mockRepo := new(MockPostRepo)
	mockSearchRepo := new(MockPostSearchRepo)
	var (
		id = uint(7)
		idFake = uint(14)
		postBF = models.Post{
			ID: 7,
			UserID: 7,
			Email: "nangcuong",
			Content: "Posted by nangcuong",
		}
		postAF = models.Post{
			ID: 7,
			UserID: 7,
			Email: "nangcuong",
			Content: "Nice testing by nangcuong",
		}
		contentAfter = "Nice testing by nangcuong"
	)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Find", id).Return(postBF, nil).Once()
		mockRepo.On("Update", postAF).Return(nil).Once()
		mockSearchRepo.On("Update", postAF).Return(nil).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		err := postService.UpdatePost(id, contentAfter)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Does-not-exist-post", func(t *testing.T) {
		mockRepo.On("Find", idFake).Return(models.Post{}, gorm.ErrRecordNotFound).Once()

		postService := NewPostService(mockRepo, mockSearchRepo)
		err := postService.UpdatePost(idFake, contentAfter)

		mockRepo.AssertExpectations(t)
		assert.Error(t, err)
	})
}
