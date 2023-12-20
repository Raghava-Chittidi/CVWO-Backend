package util

import (
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
)

func Migrate() {
	// database.DB.Migrator().CreateTable(&models.User{})
	database.DB.Migrator().CreateTable(&models.Thread{})
	// database.DB.Migrator().CreateTable(&models.Category{})
	// database.DB.Migrator().CreateTable(&models.Comment{})
	// database.DB.AutoMigrate(&models.User{}, &models.Thread{}, &models.Category{}, &models.Comment{})
}