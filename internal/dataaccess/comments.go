package data

import (
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
)

func GetAllComments() ([]*models.Comment, error) {
	var comments []*models.Comment
	result := database.DB.Table("comments").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}

// func GetAllPreloadedThreads() ([]*models.Thread, error) {
// 	var threads []*models.Thread
// 	result := database.DB.Table("threads").Preload("User", func(tx *gorm.DB) *gorm.DB {
// 		return tx.Select("id", "email", "username")
// 	}).Preload("Category").Find(&threads)

// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return threads, nil
// }

func GetCommentsByThreadId(id int) ([]*models.Comment, error) {
	var comments []*models.Comment
	result := database.DB.Table("comments").Where("thread_id = ?", id).First(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}
