package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // "-" agar tidak muncul di JSON response
	Name      string         `gorm:"not null" json:"name"`
	Posts     []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments  []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
