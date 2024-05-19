package repository

import "store.api/model"

type CardRepository interface {
	// TODO remove
	FindAll() []*model.Card
	Save(*model.Card) error
}
