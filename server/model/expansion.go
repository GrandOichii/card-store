package model

type Expansion struct {
	ID string `gorm:"not null,primaryKey" json:"id"`

	ShortName string `gorm:"not null" json:"shortName"`
	FullName  string `gorm:"not null" json:"fullName"`
}
