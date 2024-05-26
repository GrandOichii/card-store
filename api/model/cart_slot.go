package model

import "gorm.io/gorm"

type CartSlot struct {
	gorm.Model

	Amount uint `gorm:"not null" json:"amount"`

	CardID uint `gorm:"not null" json:"cardId"`
	Card   Card `json:"-"`

	CartID uint `gorm:"not null" json:"cartId"`
}
