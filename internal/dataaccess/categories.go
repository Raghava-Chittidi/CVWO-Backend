package data

import (
	"github.com/CVWO-Backend/internal/database"
	"github.com/CVWO-Backend/internal/models"
)

func GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	result := database.DB.Table("categories").Where("name = ?", name).First(&category)
	if result.Error != nil {
		return nil, result.Error
	}

	return &category, nil
}

func GetCategoryById(id int) (*models.Category, error) {
	var category models.Category
	result := database.DB.Table("categories").Where("id = ?", id).First(&category)
	if result.Error != nil {
		return nil, result.Error
	}

	return &category, nil
}

func GetAllCategories() ([]string, error) {
	var categories []string
	result := database.DB.Table("categories").Select("name").Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

