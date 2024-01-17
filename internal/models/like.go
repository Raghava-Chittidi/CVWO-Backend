package models

import "gorm.io/gorm"

type Favourite struct {
	gorm.Model
	UserID int `json:"userId"`
	User User `json:"user"`
	ThreadID int `json:"threadId"`
}

type ThreadLike struct {
	gorm.Model
	UserID int `json:"userId"`
	User User `json:"user"`
	ThreadID int `json:"threadId"`
}

type CommentLike struct {
	gorm.Model
	UserID int `json:"userId"`
	User User `json:"user"`
	CommentID int `json:"commentId"`
}