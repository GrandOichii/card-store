package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username     string `gorm:"not null;unique" json:"username"`
	PasswordHash string `gorm:"not null" json:"passwordHash"`
	Email        string `gorm:"not null" json:"email"`
	IsAdmin      bool   `gorm:"not null" json:"isAdmin"`

	Verified bool `gorm:"not null" json:"verified"`
}
