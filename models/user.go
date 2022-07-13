package models

import (
	"time"
)

type User struct {
	ID		uint `gorm:"primaryKey" json:"id"`
	Name 	string `json:"name" form:"name"`
	Email	string `gorm:"unique" json:"email" form:"email"`
	Password	string `json:"password" form:"password"`
	Role	string `json:"role"`
	Create_At	time.Time `json:"create_at"`
}

type Authen struct {
	Email 	string `json:"email"`
	Password	string `json:"password"`
	Captcha		string `json:"captcha"`
}

type ChangePassRequest struct {
	Email	string `json:"email"`
	OldPassword 	string `json:"oldPassword"`
	NewPassword 	string `json:"newPassword"`
	ConfirmPassword 	string `json:"confirmPassword"`
}

type PasswordResetToken struct {
	Email	string `gorm:"primaryKey" json:"email"`
	Token	[]byte `json:"token" gorm:"primaryKey unique"`
	ExpriedIn	time.Time `json:"expriedIn"`
}

type UserRepo interface {
	Select() ([]User, error)
	Find(email string) (User, error)
	Create(user User) (error)
	Insert(user User) (error)
	Update(user User) (error)
	Delete(email string) (error)
}

type ResetPasswordRepo interface {
	FindUser(hashToken []byte) (PasswordResetToken, error)
	FindAllToken(email string) ([]PasswordResetToken, error)
	Create(PasswordResetToken) (error)
	DeleteAll(email string) (error)
}