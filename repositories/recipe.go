package repositories

import (
	"fmt"

	"app/models"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	BeginTransaction() *gorm.DB
	CreateRecipe(tx *gorm.DB, recipe *models.Recipe) error
	CreateRecipeIngredient(tx *gorm.DB, ingredient *models.RecipeIngredient) error
	CreateRecipeStep(tx *gorm.DB, step *models.RecipeStep) error
	CreateRecipeTag(tx *gorm.DB, tag *models.RecipeTag) error
	UpdateRecipe(tx *gorm.DB, recipe *models.Recipe, ingredients []models.RecipeIngredient, steps []models.RecipeStep, tags []models.RecipeTag) error
}

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *recipeRepository) CreateRecipe(tx *gorm.DB, recipe *models.Recipe) error {
	return tx.Omit("search_vector").Create(recipe).Error
}

func (r *recipeRepository) CreateRecipeIngredient(tx *gorm.DB, ingredient *models.RecipeIngredient) error {
	return tx.Create(ingredient).Error
}

func (r *recipeRepository) CreateRecipeStep(tx *gorm.DB, step *models.RecipeStep) error {
	return tx.Create(step).Error
}

func (r *recipeRepository) CreateRecipeTag(tx *gorm.DB, tag *models.RecipeTag) error {
	return tx.Create(tag).Error
}

func (r *recipeRepository) UpdateRecipe(tx *gorm.DB, recipe *models.Recipe, ingredients []models.RecipeIngredient, steps []models.RecipeStep, tags []models.RecipeTag) error {
	// Verify recipe exists and is owned by the user
	var existingRecipe models.Recipe
	if err := tx.Where("id = ? AND creator_id = ?", recipe.ID, recipe.CreatorID).First(&existingRecipe).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("recipe not found or not owned by user")
		}
		return fmt.Errorf("failed to verify recipe: %w", err)
	}

	// Update recipe
	if err := tx.Omit("search_vector").Updates(recipe).Error; err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}

	// Delete existing associated data
	if err := tx.Where("recipe_id = ?", recipe.ID).Delete(&models.RecipeIngredient{}).Error; err != nil {
		return fmt.Errorf("failed to delete recipe_ingredient: %w", err)
	}

	if err := tx.Where("recipe_id = ?", recipe.ID).Delete(&models.RecipeStep{}).Error; err != nil {
		return fmt.Errorf("failed to delete recipe_step: %w", err)
	}

	if err := tx.Where("recipe_id = ?", recipe.ID).Delete(&models.RecipeTag{}).Error; err != nil {
		return fmt.Errorf("failed to delete recipe_tag: %w", err)
	}

	// Insert new associated data
	for _, ing := range ingredients {
		if err := r.CreateRecipeIngredient(tx, &ing); err != nil {
			return fmt.Errorf("failed to create recipe ingredient: %w", err)
		}
	}

	for _, step := range steps {
		if err := r.CreateRecipeStep(tx, &step); err != nil {
			return fmt.Errorf("failed to create recipe step: %w", err)
		}
	}

	for _, tag := range tags {
		if err := r.CreateRecipeTag(tx, &tag); err != nil {
			return fmt.Errorf("failed to create recipe tag: %w", err)
		}
	}

	return nil
}
