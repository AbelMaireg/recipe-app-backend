package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"app/framework"
	"app/models"
	"app/services"
	"app/utils"

	"github.com/google/uuid"
)

type CreateRecipeInput struct {
	Title           string                  `json:"title"`
	CategoryID      string                  `json:"category_id"`
	PreparationTime int64                   `json:"preparation_time"`
	Ingredients     []RecipeIngredientInput `json:"ingredients"`
	Steps           []RecipeStepInput       `json:"steps"`
	Tags            []RecipeTagInput        `json:"tags"`
}

type RecipeIngredientInput struct {
	IngredientID string  `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}

type RecipeStepInput struct {
	Index       int    `json:"index"`
	Description string `json:"description"`
}

type RecipeTagInput struct {
	TagID string `json:"tag_id"`
}

type CreateRecipeInputWrapper struct {
	Arg1 CreateRecipeInput `json:"arg1"`
}

type CreateRecipeResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatorID string `json:"creator_id"`
	CreatedAt string `json:"created_at"`
}

type CreateRecipeHandler struct {
	recipeService services.RecipeService
}

func (h *CreateRecipeHandler) Handle(w http.ResponseWriter, r *http.Request, action framework.HasuraAction) {
	var wrapper CreateRecipeInputWrapper
	if err := json.Unmarshal(action.Input, &wrapper); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid input format: "+err.Error())
		return
	}

	input := wrapper.Arg1
	if input.Title == "" || input.CategoryID == "" || len(input.Steps) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "MISSING_REQUIRED_FIELDS", "Title, category_id, and at least one step are required")
		return
	}
	if _, err := uuid.Parse(input.CategoryID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_UUID", "Invalid category_id: must be a valid UUID")
		return
	}
	for i, ing := range input.Ingredients {
		if _, err := uuid.Parse(ing.IngredientID); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_UUID", "Invalid ingredient_id at index "+strconv.Itoa(i)+": must be a valid UUID")
			return
		}
		if ing.Quantity <= 0 || ing.Unit == "" {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_INGREDIENT", "Invalid quantity or unit at ingredient index "+strconv.Itoa(i))
			return
		}
	}
	// Validate unique step indices
	stepIndices := make(map[int]bool)
	for i, step := range input.Steps {
		if step.Index < 1 || step.Description == "" {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_STEP", "Invalid step index or description at step index "+strconv.Itoa(i))
			return
		}
		if stepIndices[step.Index] {
			utils.WriteError(w, http.StatusBadRequest, "DUPLICATE_STEP_INDEX", "Duplicate step index "+strconv.Itoa(step.Index)+" provided")
			return
		}
		stepIndices[step.Index] = true
	}
	for i, tag := range input.Tags {
		if _, err := uuid.Parse(tag.TagID); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_UUID", "Invalid tag_id at index "+strconv.Itoa(i)+": must be a valid UUID")
			return
		}
	}

	userID, ok := action.SessionVariables["x-hasura-user-id"]
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized: Missing user ID in session")
		return
	}

	recipe := &models.Recipe{
		ID:              uuid.NewString(),
		Title:           input.Title,
		CategoryID:      input.CategoryID,
		CreatorID:       userID,
		PreparationTime: input.PreparationTime,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	var recipeIngredients []models.RecipeIngredient
	for _, ing := range input.Ingredients {
		recipeIngredients = append(recipeIngredients, models.RecipeIngredient{
			RecipeID:     recipe.ID,
			IngredientID: ing.IngredientID,
			Quantity:     ing.Quantity,
			Unit:         ing.Unit,
			CreatedAt:    time.Now(),
		})
	}

	var recipeSteps []models.RecipeStep
	for _, step := range input.Steps {
		recipeSteps = append(recipeSteps, models.RecipeStep{
			ID:          uuid.NewString(),
			RecipeID:    recipe.ID,
			Index:       step.Index,
			Description: step.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	var recipeTags []models.RecipeTag
	for _, tag := range input.Tags {
		recipeTags = append(recipeTags, models.RecipeTag{
			RecipeID:  recipe.ID,
			TagID:     tag.TagID,
			CreatedAt: time.Now(),
		})
	}

	if err := h.recipeService.CreateRecipe(userID, recipe, recipeIngredients, recipeSteps, recipeTags); err != nil {
		if strings.Contains(err.Error(), "foreign key") || strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_REFERENCE", "Invalid reference: Category, ingredient, or tag does not exist")
		} else if strings.Contains(err.Error(), "creator_id") {
			utils.WriteError(w, http.StatusForbidden, "FORBIDDEN_CREATOR", "Creator ID does not match authenticated user")
		} else if strings.Contains(err.Error(), "search_vector") {
			utils.WriteError(w, http.StatusInternalServerError, "SEARCH_VECTOR_ERROR", "Failed to process full-text search data")
		} else if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "23505") {
			utils.WriteError(w, http.StatusBadRequest, "DUPLICATE_STEP_INDEX", "Duplicate step index provided")
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to create recipe: "+err.Error())
		}
		return
	}

	response := CreateRecipeResponse{
		ID:        recipe.ID,
		Title:     recipe.Title,
		CreatorID: recipe.CreatorID,
		CreatedAt: recipe.CreatedAt.Format(time.RFC3339),
	}

	utils.EncodeJSON(w, response)
}

func RegisterCreateRecipeHandler(recipeService services.RecipeService) {
	dispatcher := framework.GetActionDispatcher(&DefaultHandler{})
	dispatcher.RegisterHandler("createRecipe", &CreateRecipeHandler{recipeService: recipeService})
}
