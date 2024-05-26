package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	UserID uint       `gorm:"not null" json:"userId"`
	Cards  []CartSlot `json:"cards"`
}
