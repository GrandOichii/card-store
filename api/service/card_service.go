package service

import "store.api/dto"

type CardService interface {
	// TODO remove
	GetAll() []*dto.GetCard
	Add(*dto.CreateCard, string) (*dto.GetCard, error)
	GetById(id uint) (*dto.GetCard, error)
}
