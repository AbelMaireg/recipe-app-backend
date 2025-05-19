package services

import (
	"fmt"

	"app/models"
	"app/repositories"
)

type RecipeService interface {
	CreateRecipe(userID string, recipe *models.Recipe, ingredients []models.RecipeIngredient, steps []models.RecipeStep, tags []models.RecipeTag) error
	UpdateRecipe(userID string, recipe *models.Recipe, ingredients []models.RecipeIngredient, steps []models.RecipeStep, tags []models.RecipeTag) error
}

type recipeService struct {
	recipeRepository repositories.RecipeRepository
}

func NewRecipeService(recipeRepo repositories.RecipeRepository) RecipeService {
	return &recipeService{recipeRepository: recipeRepo}
}

func (s *recipeService) CreateRecipe(userID string, recipe *models.Recipe, ingredients []models.RecipeIngredient, steps []models.RecipeStep, tags []models.RecipeTag) error {
	if recipe.CreatorID != userID {
		return fmt.Errorf("creator_id must match authenticated user")
	}

	tx := s.recipeRepository.BeginTransaction()
	if err := s.recipeRepository.CreateRecipe(tx, recipe); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create recipe: %w", err)
	}
	for _, ing := range ingredients {
		if err := s.recipeRepository.CreateRecipeIngredient(tx, &ing); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create recipe ingredient: %w", err)
		}
	}
	for _, step := range steps {
		if err := s.recipeRepository.CreateRecipeStep(tx, &step); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create recipe step: %w", err)
		}
	}
	for _, tag := range tags {
		if err := s.recipeRepository.CreateRecipeTag(tx, &tag); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create recipe tag: %w", err)
		}
	}
	tx.Commit()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *recipeService) UpdateRecipe(userID string, recipe *models.Recipe, ingredients []models.RecipeIngredient, steps []models.RecipeStep, tags []models.RecipeTag) error {
	if recipe.CreatorID != userID {
		return fmt.Errorf("creator_id must match authenticated user")
	}

	tx := s.recipeRepository.BeginTransaction()
	if err := s.recipeRepository.UpdateRecipe(tx, recipe, ingredients, steps, tags); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update recipe: %w", err)
	}

	tx.Commit()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
return nil
}
