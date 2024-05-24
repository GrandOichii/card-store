package service

import (
	"store.api/dto"
	"store.api/query"
)

type CardService interface {
	Add(*dto.CreateCard, uint) (*dto.GetCard, error)
	GetById(id uint) (*dto.GetCard, error)
	Query(query *query.CardQuery) []*dto.GetCard
}
