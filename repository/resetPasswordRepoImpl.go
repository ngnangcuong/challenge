package repository

import (
	"challenge3/models"
	"github.com/jinzhu/gorm"
)

type resetPassRepoImpl struct {
	DB	*gorm.DB
}

func NewResetPasswordRepo(db *gorm.DB) models.ResetPasswordRepo {
	return &resetPassRepoImpl{
		DB: db,
	}
}

func (r *resetPassRepoImpl) FindUser(hashToken []byte) (models.PasswordResetToken, error) {
	var doc models.PasswordResetToken
	result := r.DB.Where("token = ?", string(hashToken)).First(&doc)
	if result.Error != nil {
		return models.PasswordResetToken{}, result.Error
	}
	return doc, nil
}

func (r *resetPassRepoImpl) FindAllToken(email string) ([]models.PasswordResetToken, error) {
	var docs []models.PasswordResetToken
	result := r.DB.Where("email = ?", email).Find(&docs)
	if result.Error != nil {
		return nil, result.Error
	}

	return docs, nil
}

func (r *resetPassRepoImpl) Create(newDocument models.PasswordResetToken) error {
	result := r.DB.Create(&newDocument)
	return result.Error
}

func (r *resetPassRepoImpl) DeleteAll(email string) error {
	//result := r.DB.Where("email = ?").Delete(&models.PasswordResetToken{})
	result := r.DB.Exec("DELETE from password_reset_tokens where email = '" + email + "'")
	return result.Error
}