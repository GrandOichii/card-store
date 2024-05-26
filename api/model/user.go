package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username     string `gorm:"not null;unique"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"not null"`
	IsAdmin      bool   `gorm:"not null"`

	Verified bool `gorm:"not null"`

	Cart Cart
}
