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
	result := database.DB.Table("threads").Preload("User").
			  Preload("Comments", func(tx *gorm.DB) *gorm.DB {
				return tx.Order("created_at DESC")
			  }).Preload("Comments.User").Preload("Category").Order("created_at DESC, title").Find(&threads)

	if result.Error != nil {
		return nil, result.Error
	}

	return threads, nil
}

func GetPreloadedThreadById(id int) (*models.Thread, error) {
	var thread models.Thread
	result := database.DB.Table("threads").Preload("User").Preload("Comments.User").
						  Preload("Category").Where("id = ?", id).First(&thread)
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
