package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Email string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Threads []Thread `json:"threads"`
	Comments []Comment `json:"comments"`
}

