package model

import "gorm.io/gorm"

type Collection struct {
	gorm.Model

	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	Cards []CollectionSlot `json:"cards"`

	OwnerID uint `gorm:"not null" json:"ownerId"`
}
