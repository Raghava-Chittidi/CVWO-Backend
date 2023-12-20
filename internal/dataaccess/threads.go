package data

import (
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
	"gorm.io/gorm"
)

func GetAllThreads() ([]*models.Thread, error) {
	var threads []*models.Thread
	result := database.DB.Table("threads").Find(&threads)
	if result.Error != nil {
		return nil, result.Error
	}

	return threads, nil
}

func GetAllPreloadedThreads() ([]*models.Thread, error) {
	var threads []*models.Thread
	result := database.DB.Table("threads").Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "email", "username")
	}).Preload("Category").Find(&threads)

	if result.Error != nil {
		return nil, result.Error
	}

	return threads, nil
}

func GetThreadById(id int) (*models.Thread, error) {
	var thread models.Thread
	result := database.DB.Table("threads").Where("id = ?", id).First(&thread)
	if result.Error != nil {
		return nil, result.Error
	}

	return &thread, nil
}

// func GetThreadsByUsername(username string) (*models.User, error) {
// 	var user models.User
// 	result := database.DB.Table("users").Where("username = ?", username).First(&user)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return &user, nil
// }
