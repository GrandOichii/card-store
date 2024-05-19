package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username     string `json:"username" gorm:"not null;unique"`
	PasswordHash string `json:"passHash" gorm:"not null"`
	Email        string `json:"email" gorm:"not null"`

	Verified bool `json:"verified" gorm:"not null"`
}
