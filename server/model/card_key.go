package model

type CardKey struct {
	ID      string `gorm:"not null,primarKey" json:"id"`
	EngName string `gorm:"not null" json:"-"`
}
