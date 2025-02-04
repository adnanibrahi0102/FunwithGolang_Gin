package models

import "time"

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"unique"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"` // "-" to hide password from json response
	CreatedAt time.Time
	UpdatedAt time.Time
}
