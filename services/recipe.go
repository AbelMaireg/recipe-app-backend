package services

import (
	"fmt"

	"app/models"
	"app/repositories"
	"github.com/google/uuid"
)

type RecipeService interface {
	SaveRecipePicture(recipeID uuid.UUID, path string) (*models.RecipePicture, error)
	FindRecipePictureByID(id uuid.UUID) (*models.RecipePicture, error)
}

type recipeService struct {
	repository repositories.RecipeRepository
}

func NewRecipeService(repository repositories.RecipeRepository) RecipeService {
	return &recipeService{repository: repository}
}

func (r *recipeService) SaveRecipePicture(recipeID uuid.UUID, path string) (*models.RecipePicture, error) {
	picture := &models.RecipePicture{
		ID:       uuid.New(),
		RecipeId: recipeID,
		Path:     path,
	}

	if err := r.repository.SaveRecipePicture(*picture); err != nil {
		return nil, fmt.Errorf("failed to save recipe picture: %w", err)
	}

	return picture, nil
}

func (r *recipeService) FindRecipePictureByID(id uuid.UUID) (*models.RecipePicture, error) {
	picture, err := r.repository.FindRecipePictureByID(id.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find recipe picture: %w", err)
	}
	return picture, nil
}
