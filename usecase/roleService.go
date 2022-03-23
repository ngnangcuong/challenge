package usecase

import (
	"challenge3/models"
	"fmt"
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
	roleCheck, _ := r.Find(name)
	if roleCheck.Name != "" {
		return fmt.Errorf("Existed Role")
	}

	var role = models.Role{
		Name: name,
		Permission: permission,
	}

	return r.repo.Create(role)
}	