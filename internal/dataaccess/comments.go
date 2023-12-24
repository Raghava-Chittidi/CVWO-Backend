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

func GetCommentsByThreadId(id int) ([]*models.Comment, error) {
	var comments []*models.Comment
	result := database.DB.Table("comments").Where("thread_id = ?", id).Order("created_at DESC").First(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}

func GetPreloadedCommentById(id int) (*models.Comment, error) {
	var comment models.Comment
	result := database.DB.Table("comments").Preload("User").Where("id = ?", id).First(&comment)
	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func DeleteCommentById(id int) (error) {
	result := database.DB.Table("comments").Delete(&models.Comment{}, id)
	return result.Error
}
