package usecase

import (
	"challenge3/models"
)

type CheckPassService struct {
	repo models.CheckPassRepo 
}

func NewCheckPassService(repo models.CheckPassRepo) *CheckPassService {
	return &CheckPassService{
		repo: repo,
	}
}

func (c *CheckPassService) FindCheck(ip string, email string) (*models.CheckPass, error) {
	return c.repo.Find(ip, email)
}

func (c *CheckPassService) UpdateCheck(ip string, email string) error {
	_, err := c.repo.Find(ip, email)
	if err != nil {
		newCheckPass := models.CheckPass{
			IpAddress: ip,
			Email: email,
			FailedLogin: 1,
		}
		_, err = c.repo.Init(&newCheckPass)
		
		return err
	}

	return c.repo.Update(ip, email)
}

func (c *CheckPassService) DeleteCheck(ip string, email string) error {
	checkPassDoc, err := c.repo.Find(ip, email)
	if err != nil {
		return err
	}

	return c.repo.Delete(checkPassDoc)
}