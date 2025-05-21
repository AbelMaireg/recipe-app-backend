package repositories

import (
	"app/models"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	SaveRecipePicture(picture models.RecipePicture) error
	FindRecipePictureByID(id string) (*models.RecipePicture, error)
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

func (r *recipeRepository) FindRecipePictureByID(id string) (*models.RecipePicture, error) {
	var picture models.RecipePicture
	err := r.db.Joins("JOIN recipe ON recipe.id = recipe_picture.recipe_id").
		Where("recipe_picture.id = ?", id).
		First(&picture).Error
	if err != nil {
		return nil, err
	}
	return &picture, nil
}
