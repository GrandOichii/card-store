package model

import "gorm.io/gorm"

type Collection struct {
	gorm.Model

	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`

	Cards []CardSlot

	OwnerID uint `gorm:"not null"`
	Owner   User
}
