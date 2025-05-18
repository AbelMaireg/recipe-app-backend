package models

import (
	"time"

	"gorm.io/gorm"
)

type Recipe struct {
	ID              string         `gorm:"type:uuid;primaryKey"`
	Title           string         `gorm:"type:varchar(255);not null"`
	CategoryID      string         `gorm:"type:uuid;not null"`
	CreatorID       string         `gorm:"type:uuid;not null"`
	PreparationTime int64          `gorm:"not null"`
	LikeCount       int64          `gorm:"default:0"`
	RatingCount     int64          `gorm:"default:0"`
	AverageRating   float64        `gorm:"type:decimal(3,2);default:0.00"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;not null"`
	UpdatedAt       time.Time      `gorm:"type:timestamptz;not null"`
	DeletedAt       gorm.DeletedAt `gorm:"type:timestamptz"`
	SearchVector    string         `gorm:"type:tsvector"`
}

func (Recipe) TableName() string {
	return "recipe"
}

type RecipeIngredient struct {
	RecipeID     string    `gorm:"type:uuid;primaryKey"`
	IngredientID string    `gorm:"type:uuid;primaryKey"`
	Quantity     float64   `gorm:"type:decimal;not null"`
	Unit         string    `gorm:"type:varchar(50);not null"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null"`
}

func (RecipeIngredient) TableName() string {
	return "recipe_ingredient"
}

type RecipeStep struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	RecipeID    string    `gorm:"type:uuid;not null"`
	Index       int       `gorm:"not null"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null"`
}

func (RecipeStep) TableName() string {
	return "recipe_step"
}

type RecipeTag struct {
	RecipeID  string    `gorm:"type:uuid;primaryKey"`
	TagID     string    `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null"`
}

func (RecipeTag) TableName() string {
	return "recipe_tag"
}
