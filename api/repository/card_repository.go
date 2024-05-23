package repository

import "store.api/model"

type CardRepository interface {
	Save(*model.Card) error
	FindById(id uint) *model.Card
	FindByType(cType string) []*model.Card
}
