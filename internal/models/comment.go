package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	UserID int `json:"userId"`
	User User `json:"user"`
	ThreadID int `json:"threadId"`
	Thread Thread `json:"thread"`
}