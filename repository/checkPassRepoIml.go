package repository

import (
	"challenge3/models"
	"github.com/jinzhu/gorm"
)

type checkPassRepoIml struct {
	DB *gorm.DB
}

func NewCheckPassRepo(db *gorm.DB) models.CheckPassRepo {
	return &checkPassRepoIml{
		DB: db,
	}
}

func (c *checkPassRepoIml) Init(newCheckPass *models.CheckPass) (*models.CheckPass, error) {
	result := c.DB.Create(newCheckPass)
	return newCheckPass, result.Error
}

func (c *checkPassRepoIml) Find(ip string, email string) (*models.CheckPass, error) {
	var checkPassDoc models.CheckPass
	result := c.DB.Where("ip_address = ? AND email = ?", ip, email).First(&checkPassDoc)
	if result.Error != nil {
		return &models.CheckPass{}, result.Error
	}

	return &checkPassDoc, nil
}

func (c *checkPassRepoIml) Update(ip string, email string) error {
	checkPassDoc, err := c.Find(ip, email)
	if err != nil {
		return err
	}
	result := c.DB.Model(checkPassDoc).Update("failed_login", checkPassDoc.FailedLogin + 1)
	return result.Error
}

func (c *checkPassRepoIml) Delete(checkPassDoc *models.CheckPass) error {
	result := c.DB.Delete(checkPassDoc)
	return result.Error
}	