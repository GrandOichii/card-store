package model

import "gorm.io/gorm"

type CollectionSlot struct {
	gorm.Model

	Amount uint `gorm:"not null" json:"amount"`

	CardID uint `gorm:"not null" json:"cardId"`
	Card   Card `json:"-"`

	CollectionID uint `gorm:"not null" json:"collectionId"`
}
