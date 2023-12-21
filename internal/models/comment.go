package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	UserID int `json:"userId"`
	User User `json:"user" gorm:"preload:true"`
	ThreadID int `json:"threadId"`
}