package usecase

import (
	"challenge3/models"
)

type PostService struct {
	repo models.PostRepo
}

func NewPostService(repo models.PostRepo) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (p *PostService) GetListPost() ([]models.Post, error) {
	return p.repo.Select()
}

func (p *PostService) CreatePost(post models.Post) (models.Post, error) {
	err := p.repo.Create(post)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (p *PostService) DeletePost(id uint) (error) {
	_, err := p.Find(id)
	if err != nil {
		return err
	}

	return p.repo.Delete(id)
}

func (p *PostService) UpdatePost(id uint, content string) (error) {
	postCheck, err := p.Find(uint(id))
	if err != nil {
		return err
	}

	postCheck.Content = content

	return p.repo.Update(postCheck)
}

func (p *PostService) Find(id uint) (models.Post, error) {
	return p.repo.Find(id)
}