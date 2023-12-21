package handlers

import (
	"errors"
	"net/http"

	data "github.com/CVWO-Backend/internal/dataaccess"
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
	"github.com/CVWO-Backend/internal/util"
)

// const (
// 	ListUsers = "users.HandleList"

// 	SuccessfulListUsersMessage = "Successfully listed users"
// 	ErrRetrieveDatabase        = "Failed to retrieve database in %s"
// 	ErrRetrieveUsers           = "Failed to retrieve users in %s"
// 	ErrEncodeView              = "Failed to retrieve users in %s"
// )

func Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	} {
		Status: "active",
		Message: "ForumZone running!",
		Version: "1.0.0",
	}

	_ = util.WriteJSON(w, payload, http.StatusOK)
}

func CreateThread(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title string `json:"title"`
		Content string `json:"content"`
		ImageUrl string `json:"imageUrl"`
		Username string `json:"username"`
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

	user, err := data.GetUserByUsername(payload.Username)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	newThread := &models.Thread{Title: payload.Title, Content: payload.Content, ImageUrl: payload.ImageUrl, 
								UserID: int(user.ID), CategoryID: int(category.ID)}
	result := database.DB.Table("threads").Create(newThread)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	var message = struct {
		Message string `json:"message"`
	} {
		Message: "Thread successfully created!",
	}
	
	util.WriteJSON(w, message, http.StatusCreated)
}

func GetThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := data.GetAllPreloadedThreads()
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, threads, http.StatusOK)
}
