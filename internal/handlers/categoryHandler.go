package handlers

import (
	"net/http"

	data "github.com/CVWO-Backend/internal/dataaccess"
	"github.com/CVWO-Backend/internal/util"
)

// func CreateCategory(w http.ResponseWriter, r *http.Request) {
// 	cat1 := &models.Category{Name: "Environment", Threads: []models.Thread{}}
// 	cat2 := &models.Category{Name: "Finance", Threads: []models.Thread{}}
// 	cat3 := &models.Category{Name: "Politics", Threads: []models.Thread{}}
// 	cat4 := &models.Category{Name: "General", Threads: []models.Thread{}}
// 	cat5 := &models.Category{Name: "Healthcare", Threads: []models.Thread{}}
// 	cat6 := &models.Category{Name: "Music", Threads: []models.Thread{}}
// 	cat7 := &models.Category{Name: "Movies", Threads: []models.Thread{}}
// 	cat8 := &models.Category{Name: "Education", Threads: []models.Thread{}}
// 	category := []*models.Category{cat1, cat2, cat3, cat4, cat5, cat6, cat7, cat8}

// 	result := database.DB.Table("categories").Create(category)
// 	if result.Error != nil {
// 		util.ErrorJSON(w, result.Error, http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := data.GetAllCategories()
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, categories, http.StatusOK)
}