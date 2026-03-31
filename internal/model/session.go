package model

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Title     string         `gorm:"type:varchar(100);not null" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
