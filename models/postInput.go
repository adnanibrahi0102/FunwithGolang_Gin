package models

import "time"

type PostInput struct {
	Title     string    `form:"title" binding:"required"`
	Content   string    `form:"content" binding:"required"`
	CreatedAt time.Time `form:"created_at"` // Optional, no binding tag
	UpdatedAt time.Time `form:"updated_at"` // Optional, no binding tag
}

