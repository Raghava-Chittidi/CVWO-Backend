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
	result := database.DB.Table("threads").Preload("User").Preload("Likes.User").
			  Preload("Comments", func(tx *gorm.DB) *gorm.DB {
				return tx.Order("created_at DESC")
			  }).Preload("Comments.User").Preload("Comments.Likes").Preload("Comments.Likes.User").
			  Preload("Favourites.User").Preload("Category").Order("created_at DESC, title").Find(&threads)

	if result.Error != nil {
		return nil, result.Error
	}

	return threads, nil
}

func GetPreloadedThreadById(id int) (*models.Thread, error) {
	var thread models.Thread
	result := database.DB.Table("threads").Preload("User").Preload("Comments.User").Preload("Comments.Likes").
						  Preload("Comments.Likes.User").Preload("Category").
						  Preload("Likes.User").Preload("Favourites.User").Where("id = ?", id).First(&thread)
	if result.Error != nil {
		return nil, result.Error
	}

	return &thread, nil
}

func GetThreadById(id int) (*models.Thread, error) {
	var thread models.Thread
	result := database.DB.Table("threads").Where("id = ?", id).First(&thread)
	if result.Error != nil {
		return nil, result.Error
	}

	return &thread, nil
}

func DeleteThreadById(id int) (error) {
	result := database.DB.Table("threads").Delete(&models.Thread{}, id)
	return result.Error
}
