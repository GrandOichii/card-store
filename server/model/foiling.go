package model

type Foiling struct {
	ID              string `gorm:"not null,primaryKey" json:"id"`
	Label           string `gorm:"not null" json:"label"`
	DescriptiveName string `gorm:"not null" json:"descriptiveName"`
}
