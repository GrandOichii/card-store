package model

import "gorm.io/gorm"

type CardSlot struct {
	gorm.Model

	Amount uint `gorm:"not null"`

	CardID uint `gorm:"not null"`
	Card   Card

	CollectionID uint `gorm:"not null"`
}
