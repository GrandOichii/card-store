package model

import "gorm.io/gorm"

// TODO? add card keys - different cards that are basically the same core card have the same key (example: Russian Iron Myr and English Iron Myr)

type Card struct {
	gorm.Model

	Name     string  `gorm:"not null"`
	Text     string  `gorm:"not null,type:text"`
	ImageUrl string  `gorm:""`
	Price    float32 `gorm:"not null"`

	PosterID uint `gorm:"not null"`
	Poster   User

	CardTypeID string `gorm:"not null"`
	CardType   CardType

	LanguageID string `gorm:"not null"`
	Language   Language
}
