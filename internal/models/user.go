package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Email string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	// ImageUrl string `json:"ImageUrl"`
	Threads []Thread `json:"threads"`
	Comments []Comment `json:"comments"`
}

func (u *User) VerifyPassword(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plaintext))
	if err != nil {
		return false, err
	}

	return true, nil
}

