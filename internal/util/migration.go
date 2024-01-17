package util

import (
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
)

func Migrate() {
	// database.DB.Migrator().CreateTable(&models.User{})
	// database.DB.Migrator().CreateTable(&models.Comment{})
	// database.DB.Migrator().CreateTable(&models.Thread{})
	// database.DB.Migrator().CreateTable(&models.Category{})
	database.DB.Migrator().CreateTable(&models.ThreadLike{})
	// database.DB.Migrator().CreateTable(&models.CommentLike{})
	// database.DB.Migrator().CreateTable(&models.Favourite{})
	// database.DB.AutoMigrate(&models.Comment{}, &models.Thread{}, &models.CommentLike{}, &models.ThreadLike{}, &models.Favourite{})
}