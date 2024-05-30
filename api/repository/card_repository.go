package repository

import (
	"store.api/model"
	"store.api/query"
)

type CardRepository interface {
	Save(*model.Card) error
	FindById(id uint) *model.Card
	Update(*model.Card) error
	UpdatePrice(id uint, price float32) (*model.Card, error)
	UpdateInStockAmount(id uint, amount uint) (*model.Card, error)
	Query(query *query.CardQuery) []*model.Card
	Count() int64
}
