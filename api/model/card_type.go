package model

type CardType struct {
	ID        string `gorm:"not null,primaryKey" json:"id"`
	LongName  string `gorm:"not null" json:"longName"`
	ShortName string `gorm:"not null" json:"-"`
}
