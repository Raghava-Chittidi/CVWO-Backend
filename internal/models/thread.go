package models

import (
	"gorm.io/gorm"
)

type Thread struct {
	gorm.Model
	Topic string `json:"topic"`
	Content string `json:"content"`
	UserID uint `json:"userId"`
	CategoryID uint `json:"categoryId"`
	Comments []Comment `json:"comments"`
}

