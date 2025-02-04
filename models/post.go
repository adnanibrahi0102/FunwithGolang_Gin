package models

import "time"

type Post struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	Content   string `json:"content"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

