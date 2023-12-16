package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
	Threads []Thread `json:"threads"`
}