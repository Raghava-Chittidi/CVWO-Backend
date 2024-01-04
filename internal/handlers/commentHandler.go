package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/CVWO-Backend/internal/auth"
	data "github.com/CVWO-Backend/internal/dataaccess"
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
	"github.com/CVWO-Backend/internal/util"
	"github.com/go-chi/chi/v5"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Comment string `json:"comment"`
		ThreadID int `json:"threadId"`
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if payload.Comment == "" {
		util.ErrorJSON(w, errors.New("Comment cannot be empty!"))
		return
	}

	_, claims, err := auth.Auth.VerifyToken(w, r)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	user, err := data.GetUserByUsername(claims.Username)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	newComment := &models.Comment{Content: payload.Comment, UserID: int(user.ID), ThreadID: payload.ThreadID}
	result := database.DB.Table("comments").Create(newComment)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	comment, err := data.GetPreloadedCommentById(int(newComment.ID))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	
	data := util.JSONResponse{Error: false, Message: "Commented successfully!", Data: *comment}
	util.WriteJSON(w, data, http.StatusCreated)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	threads, err := data.GetAllComments()
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, threads, http.StatusOK)
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id")) 
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var payload struct {
		Content string `json:"content"`
	}

	err = util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if payload.Content == "" {
		util.ErrorJSON(w, errors.New("Comment cannot be empty!"))
		return
	}

	comment, err := data.GetPreloadedCommentById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, claims, err := auth.Auth.VerifyToken(w, r)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	user, err := data.GetUserByUsername(claims.Username)
	if err != nil || comment.User.ID != user.ID {
		util.ErrorJSON(w, errors.New("Unauthorized!"), http.StatusUnauthorized)
		return
	}

	comment.Content = payload.Content
	database.DB.Save(&comment)
	
	data := util.JSONResponse{Error: false, Message: "Edited comment successfully!"}
	util.WriteJSON(w, data, http.StatusOK)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	comment, err := data.GetPreloadedCommentById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, claims, err := auth.Auth.VerifyToken(w, r)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	user, err := data.GetUserByUsername(claims.Username)
	if err != nil || comment.User.ID != user.ID {
		util.ErrorJSON(w, errors.New("Unauthorized!"), http.StatusUnauthorized)
		return
	}

	err = data.DeleteCommentById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	data := util.JSONResponse{Error: false, Message: "Deleted comment successfully!"}
	util.WriteJSON(w, data, http.StatusOK)
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, err = data.GetCommentById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, claims, err := auth.Auth.VerifyToken(w, r)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	user, err := data.GetUserByUsername(claims.Username)
	if err != nil {
		util.ErrorJSON(w, errors.New("Unauthorized!"), http.StatusUnauthorized)
		return
	}

	newCommentLike := models.CommentLike{UserID: int(user.ID), CommentID: id}
	result := database.DB.Create(&newCommentLike)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	data := util.JSONResponse{Error: false, Message: "Liked comment!"}
	util.WriteJSON(w, data, http.StatusOK)
}

func UnlikeComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, err = data.GetCommentById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, claims, err := auth.Auth.VerifyToken(w, r)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	user, err := data.GetUserByUsername(claims.Username)
	if err != nil {
		util.ErrorJSON(w, errors.New("Unauthorized!"), http.StatusUnauthorized)
		return
	}

	commentLike, err := data.GetCommentLikeByUserAndComment(int(user.ID), id)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = data.DeleteCommentLikeById(int(commentLike.ID))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	data := util.JSONResponse{Error: false, Message: "Unliked comment!"}
	util.WriteJSON(w, data, http.StatusOK)
}
