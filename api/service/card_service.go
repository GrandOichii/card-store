package service

import "store.api/dto"

type CardService interface {
	// TODO remove
	GetAll() []*dto.GetCard
	Add(*dto.CreateCard) (*dto.GetCard, error)
}
