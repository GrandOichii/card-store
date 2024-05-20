package model

import "gorm.io/gorm"

type Card struct {
	gorm.Model

	Name     string `gorm:"not null"`
	Text     string `gorm:"not null,type:text"`
	ImageUrl string `gorm:""`

	PosterId uint `gorm:"not null"`
	Poster   User
}
