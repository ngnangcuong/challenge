package repository

import (
	"challenge3/models"
	"github.com/jinzhu/gorm"
)

type postRepoImpl struct {
	DB *gorm.DB
}

func NewPostRepo(db *gorm.DB) models.PostRepo {
	return &postRepoImpl{
		DB: db,
	}
}

func (p *postRepoImpl) Select() ([]models.Post, error) {
	var postList []models.Post
	result := p.DB.Find(&postList)

	if result.Error != nil {
		return nil, result.Error
	}

	return postList, nil
}

func (p *postRepoImpl) Delete(id uint) (error) {
	post, _ := p.Find(id)
	result := p.DB.Delete(&post)
	return result.Error
}

func (p *postRepoImpl) Update(post models.Post)	(error) {
	result := p.DB.Save(&post)
	return result.Error
}

func (p *postRepoImpl) Create(post models.Post) (models.Post, error) {
	result := p.DB.Create(&post)
	return post, result.Error
}

func (p *postRepoImpl) Find(id uint) (models.Post, error) {
	var post models.Post
	result := p.DB.Where("id = ?", id).First(&post)
	
	if result.Error != nil {
		return models.Post{}, result.Error
	}
	return post, nil
}
