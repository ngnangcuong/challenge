package usecase

import (
	"challenge3/models"
	"strconv"
)

type PostService struct {
	repo models.PostRepo
	search models.PostSearchRepo
}

func NewPostService(repo models.PostRepo, searchRepo models.PostSearchRepo) *PostService {
	return &PostService{
		repo: repo,
		search: searchRepo,
	}
}

func (p *PostService) GetListPost() ([]models.Post, error) {
	return p.repo.Select()
}

func (p *PostService) CreatePost(post models.Post) (models.Post, error) {
	postES, err := p.repo.Create(post)
	if err != nil {
		return models.Post{}, err
	}

	err = p.search.Index(postES)
	if err != nil {
		return models.Post{}, err
	}

	return postES, nil
}

func (p *PostService) DeletePost(id uint) (error) {
	_, err := p.Find(id)
	if err != nil {
		return err
	}

	err = p.repo.Delete(id)
	if err != nil {
		return err
	}

	idES := strconv.FormatUint(uint64(id), 10)
	err = p.search.Delete(idES)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) UpdatePost(id uint, content string) (error) {
	postCheck, err := p.Find(id)
	if err != nil {
		return err
	}

	postCheck.Content = content

	err = p.repo.Update(postCheck)
	if err != nil {
		return err
	}

	err = p.search.Update(postCheck)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) Find(id uint) (models.Post, error) {
	return p.repo.Find(id)
}

func (p *PostService) SearchPosts(keyword string) ([]models.Post, error) {
	return p.search.Search(keyword)
}