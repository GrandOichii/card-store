package model

type Language struct {
	ID       string `gorm:"not null,primaryKey" json:"id"`
	LongName string `gorm:"not null" json:"longName"`
}
