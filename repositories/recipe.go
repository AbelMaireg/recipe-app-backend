package repositories

import (
	"app/models"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	SaveRecipePicture(picture models.RecipePicture) error
}

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) SaveRecipePicture(picture models.RecipePicture) error {
	return r.db.Create(&picture).Error
}
