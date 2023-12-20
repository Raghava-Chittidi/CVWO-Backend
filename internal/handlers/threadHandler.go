package handlers

import (
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
	}

	category, err := data.GetCategoryByName(payload.Category)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
	}

	user, err := data.GetUserByUsername(payload.Username)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
	}

	newThread := &models.Thread{Title: payload.Title, Content: payload.Content, ImageUrl: payload.ImageUrl, 
								UserID: int(user.ID), CategoryID: int(category.ID)}
	result := database.DB.Table("threads").Create(newThread)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
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
	}

	util.WriteJSON(w, threads, http.StatusOK)
}
