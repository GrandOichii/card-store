package model

import "gorm.io/gorm"

type Card struct {
	gorm.Model

	Name          string  `gorm:"not null" json:"name"`
	Text          string  `gorm:"not null,type:text" json:"text"`
	ImageUrl      string  `gorm:"" json:"imageUrl"`
	Price         float32 `gorm:"not null" json:"price"`
	InStockAmount uint    `gorm:"not null"`

	CardKeyID string `gorm:"not null" json:"cardKeyId"`

	PosterID uint `gorm:"not null" json:"posterId"`
	Poster   User `json:"-"`

	CardTypeID string   `gorm:"not null" json:"cardTypeId"`
	CardType   CardType `json:"cardType"`

	LanguageID string   `gorm:"not null" json:"languageId"`
	Language   Language `json:"language"`

	ExpansionID string    `gorm:"not null" json:"expansionId"`
	Expansion   Expansion `json:"expansion"`

	FoilingID *string `gorm:"" json:"foilingId"`
	Foiling   Foiling `json:"foiling"`
}
