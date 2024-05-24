package service

import (
	"store.api/dto"
)

type CardService interface {
	Add(*dto.CreateCard, uint) (*dto.GetCard, error)
	GetById(id uint) (*dto.GetCard, error)
	GetByType(cType string) []*dto.GetCard
}
