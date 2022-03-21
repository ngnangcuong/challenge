package usecase

import (
	"challenge3/models"

)

type RoleService struct {
	repo models.RoleRepo
}

func NewRoleService(repo models.RoleRepo) *RoleService {
	return &RoleService{
		repo: repo,

	}	
}

func (r *RoleService) Find(name string) (models.Role, error) {
	return r.repo.Find(name)
}

func (r *RoleService) Create(name, permission string) (error) {
	var role = models.Role{
		Name: name,
		Permission: permission,
	}

	return r.repo.Create(role)
}	