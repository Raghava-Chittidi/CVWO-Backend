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

	_, claims, err := auth.Auth.VerifyAuthorisationToken(w, r)
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

	data := util.JSONResponse{Error: false, Message: "Thread successfully created!", Data: *thread}
	util.WriteJSON(w, data, http.StatusCreated)
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
	// id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/threads/"))
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