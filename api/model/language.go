package model

type Language struct {
	ID       string `gorm:"not null,primaryKey"`
	LongName string `gorm:"not null"`
}
