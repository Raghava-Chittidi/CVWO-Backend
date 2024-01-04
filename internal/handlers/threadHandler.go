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

// const (
// 	ListUsers = "users.HandleList"

// 	SuccessfulListUsersMessage = "Successfully listed users"
// 	ErrRetrieveDatabase        = "Failed to retrieve database in %s"
// 	ErrRetrieveUsers           = "Failed to retrieve users in %s"
// 	ErrEncodeView              = "Failed to retrieve users in %s"
// )

func CreateThread(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title string `json:"title"`
		Content string `json:"content"`
		ImageUrl string `json:"imageUrl"`
		Category string `json:"category"`
	}

	err := util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if payload.Category == "" {
		util.ErrorJSON(w, errors.New("Category cannot be empty!"))
		return
	}

	if payload.Title == "" {
		util.ErrorJSON(w, errors.New("Title cannot be empty!"))
		return
	}

	if payload.Content == "" {
		util.ErrorJSON(w, errors.New("Content cannot be empty!"))
		return
	}

	category, err := data.GetCategoryByName(payload.Category)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
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

	newThread := &models.Thread{Title: payload.Title, Content: payload.Content, ImageUrl: payload.ImageUrl, 
								UserID: int(user.ID), CategoryID: int(category.ID)}
	result := database.DB.Table("threads").Create(newThread)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	thread, err := data.GetPreloadedThreadById(int(newThread.ID))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	data := util.ResponseJSON{Error: false, Message: "Thread successfully created!", Data: *thread}
	util.WriteJSON(w, data, http.StatusCreated)
}

func EditThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var payload struct {
		Title string `json:"title"`
		Content string `json:"content"`
		ImageUrl string `json:"imageUrl"`
		Category string `json:"category"`
	}

	err = util.ReadJSON(w, r, &payload)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	if payload.Category == "" {
		util.ErrorJSON(w, errors.New("Category cannot be empty!"))
		return
	}

	if payload.Title == "" {
		util.ErrorJSON(w, errors.New("Title cannot be empty!"))
		return
	}

	if payload.Content == "" {
		util.ErrorJSON(w, errors.New("Content cannot be empty!"))
		return
	}

	category, err := data.GetCategoryByName(payload.Category)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_, claims, err := auth.Auth.VerifyToken(w, r)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	thread, err := data.GetThreadById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	preloadedThread, err := data.GetPreloadedThreadById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	user, err := data.GetUserByUsername(claims.Username)
	if err != nil || preloadedThread.User.ID != user.ID {
		util.ErrorJSON(w, errors.New("Unauthorized!"), http.StatusUnauthorized)
		return
	}

	thread.Title = payload.Title
	thread.ImageUrl = payload.ImageUrl
	thread.Content = payload.Content
	thread.CategoryID = int(category.ID)
	database.DB.Save(&thread)

	editedThread, err := data.GetPreloadedThreadById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	data := util.ResponseJSON{Error: false, Message: "Thread successfully edited!", Data: editedThread}
	util.WriteJSON(w, data, http.StatusOK)
}

func GetThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := data.GetAllPreloadedThreads()
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, threads, http.StatusOK)
}

func GetThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	thread, err := data.GetPreloadedThreadById(id)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, thread, http.StatusOK)
}

func DeleteThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	thread, err := data.GetPreloadedThreadById(id)
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
	if err != nil || thread.User.ID != user.ID {
		util.ErrorJSON(w, errors.New("Unauthorized!"), http.StatusUnauthorized)
		return
	}

	err = data.DeleteThreadById(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	data := util.ResponseJSON{Error: false, Message: "Deleted thread successfully!"}
	util.WriteJSON(w, data, http.StatusOK)
}

func LikeThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, err = data.GetThreadById(id)
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

	newThreadLike := models.ThreadLike{UserID: int(user.ID), ThreadID: id}
	result := database.DB.Create(&newThreadLike)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	data := util.ResponseJSON{Error: false, Message: "Liked thread!"}
	util.WriteJSON(w, data, http.StatusOK)
}

func UnlikeThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, err = data.GetThreadById(id)
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

	threadLike, err := data.GetThreadLikeByUserAndThread(int(user.ID), id)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = data.DeleteThreadLikeById(int(threadLike.ID))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	data := util.ResponseJSON{Error: false, Message: "Unliked thread!"}
	util.WriteJSON(w, data, http.StatusOK)
}

func FavouriteThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, err = data.GetThreadById(id)
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

	newFavourite := models.Favourite{UserID: int(user.ID), ThreadID: id}
	result := database.DB.Create(&newFavourite)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	data := util.ResponseJSON{Error: false, Message: "Favourited thread!"}
	util.WriteJSON(w, data, http.StatusOK)
}

func UnfavouriteThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	_, err = data.GetThreadById(id)
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

	threadLike, err := data.GetFavouriteByUserAndThread(int(user.ID), id)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = data.DeleteFavouriteById(int(threadLike.ID))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	data := util.ResponseJSON{Error: false, Message: "Removed thread from favourites!"}
	util.WriteJSON(w, data, http.StatusOK)
}