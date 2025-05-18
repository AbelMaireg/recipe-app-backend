package repositories

import (
	"app/models"

	"gorm.io/gorm"
)

type RecipeRepository interface {
	BeginTransaction() *gorm.DB
	CreateRecipe(tx *gorm.DB, recipe *models.Recipe) error
	CreateRecipeIngredient(tx *gorm.DB, ingredient *models.RecipeIngredient) error
	CreateRecipeStep(tx *gorm.DB, step *models.RecipeStep) error
	CreateRecipeTag(tx *gorm.DB, tag *models.RecipeTag) error
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
