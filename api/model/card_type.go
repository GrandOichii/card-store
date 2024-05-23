package model

type CardType struct {
	ID       string `gorm:"not null,primaryKey"`
	LongName string `gorm:"not null"`
}
