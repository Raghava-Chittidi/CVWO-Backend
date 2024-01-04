package models

import (
	"gorm.io/gorm"
)

type Thread struct {
	gorm.Model
	Title string `json:"title"`
	Content string `json:"content"`
	ImageUrl string `json:"imageUrl"`
	UserID int `json:"userId"`
	User User `json:"user"`
	CategoryID int `json:"categoryId"`
	Category Category `json:"category"`
	Comments []Comment `json:"comments" gorm:"preload:true"`
	Likes []ThreadLike `json:"likes"`
	Favourites []Favourite `json:"favourites"`
}

