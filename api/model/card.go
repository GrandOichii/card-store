package model

import "gorm.io/gorm"

type Card struct {
	gorm.Model

	Name string `gorm:"not null"`
	Text string `gorm:"not null,type:text"`

	PosterId uint `gorm:"not null"`
	Poster   User
}
