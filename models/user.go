package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Username  string         `gorm:"type:varchar(255);unique;not null"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Bio       string         `gorm:"type:text"`
	Password  string         `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index"`
}

func (User) TableName() string {
	return "user"
}
