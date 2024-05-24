package repository

import (
	"gorm.io/gorm"
	"store.api/model"
)

type CardRepository interface {
	Save(*model.Card) error
	FindById(id uint) *model.Card
	Query(applyQueryF func(*gorm.DB) *gorm.DB) []*model.Card
}
