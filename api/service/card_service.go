package service

import (
	"errors"

	"store.api/dto"
	"store.api/query"
)

var (
	ErrCardNotFound = errors.New("card not found")
)

type CardService interface {
	Add(*dto.PostCard, uint) (*dto.GetCard, error)
	GetById(id uint) (*dto.GetCard, error)
	Query(query *query.CardQuery) []*dto.GetCard
	Update(*dto.PostCard, uint) (*dto.GetCard, error)
}
