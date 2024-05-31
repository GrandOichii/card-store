package service

import (
	"errors"

	"store.api/dto"
	"store.api/query"
)

var (
	ErrCardNotFound = errors.New("card not found")
)

type CardQueryResult struct {
	Cards      []*dto.GetCard `json:"cards"`
	TotalCount int64          `json:"totalCards"`
	PerPage    uint           `json:"perPage"`
}

type CardService interface {
	Add(*dto.PostCard, uint) (*dto.GetCard, error)
	GetById(id uint) (*dto.GetCard, error)
	Query(query *query.CardQuery) *CardQueryResult
	Update(*dto.PostCard, uint) (*dto.GetCard, error)
	UpdatePrice(uint, *dto.PriceUpdate) (*dto.GetCard, error)
	UpdateInStockAmount(uint, *dto.StockedAmountUpdate) (*dto.GetCard, error)
}
