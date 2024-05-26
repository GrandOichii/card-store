package model

import "gorm.io/gorm"

// TODO? add card keys - different cards that are basically the same core card have the same key (example: Russian Iron Myr and English Iron Myr)

type Card struct {
	gorm.Model

	Name     string  `gorm:"not null" json:"name"`
	Text     string  `gorm:"not null,type:text" json:"text"`
	ImageUrl string  `gorm:"" json:"imageUrl"`
	Price    float32 `gorm:"not null" json:"price"`

	PosterID uint `gorm:"not null" json:"posterId"`
	Poster   User `json:"poster"`

	CardTypeID string   `gorm:"not null" json:"cardTypeId"`
	CardType   CardType `json:"cardType"`

	LanguageID string   `gorm:"not null" json:"languageId"`
	Language   Language `json:"language"`
}
