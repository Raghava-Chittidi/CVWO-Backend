package handlers

import (
	"errors"
	"net/http"

	data "github.com/CVWO-Backend/internal/dataaccess"
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
	"github.com/CVWO-Backend/internal/util"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
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

	user, err := data.GetUserByUsername(payload.Username)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	newComment := &models.Comment{Content: payload.Comment, UserID: int(user.ID), ThreadID: payload.ThreadID}
	result := database.DB.Table("comments").Create(newComment)
	if result.Error != nil {
		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
		return
	}

	var message = struct {
		Message string `json:"message"`
	} {
		Message: "Commented successfully!",
	}
	
	util.WriteJSON(w, message, http.StatusCreated)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	threads, err := data.GetAllPreloadedThreads()
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, threads, http.StatusOK)
}
