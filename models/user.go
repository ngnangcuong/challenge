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
}

type UserRepo interface {
	Select() ([]User, error)
	Find(email string) (User, error)
	Create(user User) (error)
	Insert(user User) (error)
	Update(user User) (error)
	Delete(email string) (error)
}