package model

type Expansion struct {
	ID string `gorm:"not null,primaryKey" json:"id"`

	ShortName string `gorm:"not null"`
	FullName  string `gorm:"not null"`
}
