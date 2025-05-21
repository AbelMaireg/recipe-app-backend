package models

import (
	"time"

	"github.com/google/uuid"
)

type RecipePicture struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	RecipeId  uuid.UUID `gorm:"type:uuid;not null"`
	Path      string    `gorm:"type:varchar(225);not null;unique"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
}

func (RecipePicture) TableName() string {
	return "recipe_picture"
}
