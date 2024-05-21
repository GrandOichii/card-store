package model

import "gorm.io/gorm"

type Card struct {
	gorm.Model

	Name     string  `gorm:"not null"`
	Text     string  `gorm:"not null,type:text"`
	ImageUrl string  `gorm:""`
	Price    float32 `gorm:"not null"`

	PosterID uint `gorm:"not null"`
	Poster   User
}
