package data

import (
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
)

func GetThreadLikeByUserAndThread(userId int, threadId int) (*models.ThreadLike, error) {
	var threadLike models.ThreadLike
	result := database.DB.Table("thread_likes").Where("user_id = ?", userId).Where("thread_id = ?", threadId).First(&threadLike)
	if result.Error != nil {
		return nil, result.Error
	}

	return &threadLike, nil
}

func DeleteThreadLikeById(id int) (error) {
	result := database.DB.Table("thread_likes").Delete(&models.ThreadLike{}, id)
	return result.Error
}

func GetCommentLikeByUserAndComment(userId int, commentId int) (*models.CommentLike, error) {
	var commentLike models.CommentLike
	result := database.DB.Table("comment_likes").Where("user_id = ?", userId).Where("comment_id = ?", commentId).First(&commentLike)
	if result.Error != nil {
		return nil, result.Error
	}

	return &commentLike, nil
}

func DeleteCommentLikeById(id int) (error) {
	result := database.DB.Table("comment_likes").Delete(&models.CommentLike{}, id)
	return result.Error
}

func GetFavouriteByUserAndThread(userId int, threadId int) (*models.Favourite, error) {
	var favourite models.Favourite
	result := database.DB.Table("favourites").Where("user_id = ?", userId).Where("thread_id = ?", threadId).First(&favourite)
	if result.Error != nil {
		return nil, result.Error
	}

	return &favourite, nil
}

func DeleteFavouriteById(id int) (error) {
	result := database.DB.Table("favourites").Delete(&models.Favourite{}, id)
	return result.Error
}