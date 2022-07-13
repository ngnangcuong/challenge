package models

import (

)

type CheckPass struct {
	IpAddress 	string `json:"ip_address"`
	Email	string `json:"email"`
	FailedLogin		int `json:"failed"`
}

type CheckPassRepo interface {
	Find(ip string, email string) (*CheckPass, error)
	Init(checkPassDoc *CheckPass) (*CheckPass, error)
	Update(ip string, email string) (error)
	Delete(checkPassDoc *CheckPass) (error)
}