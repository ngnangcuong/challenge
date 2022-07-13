package usecase

import (
	"challenge3/models"
	"crypto/sha1"
	"math/rand"
	"time"
)

type ResetPasswordService struct {
	repo models.ResetPasswordRepo
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKMNOPQRSTUVWXYZ")

func NewResetPasswordService(repo models.ResetPasswordRepo) *ResetPasswordService {
	return &ResetPasswordService{
		repo: repo,
	}
}

func (r *ResetPasswordService) GenerateResetPasswordToken(n int) string {
	tmp := make([]rune, n)
	for i := range tmp {
		tmp[i] = letters[rand.Intn(len(letters))]
	}

	return string(tmp)
}

func (r *ResetPasswordService) StoreToken(token []byte, email string) error {
	hasher := sha1.New()
	hasher.Write(token)
	result := hasher.Sum(nil)
	newDocument := models.PasswordResetToken{
		Email: email,
		Token: result,
		ExpriedIn: time.Now().Add(time.Minute * 5),
	}
	return r.repo.Create(newDocument)
}

func (r *ResetPasswordService) FindUser(token string) (models.PasswordResetToken, error) {
	hasher := sha1.New()
	hasher.Write([]byte(token))
	hashToken := hasher.Sum(nil)
	return r.repo.FindUser(hashToken)
}

func (r *ResetPasswordService) FindAllToken(email string)([]models.PasswordResetToken, error) {
	return r.repo.FindAllToken(email)
}

func (r *ResetPasswordService) DeleteAll(email string) error {
	return r.repo.DeleteAll(email)
}
