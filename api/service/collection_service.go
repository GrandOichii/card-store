package service

import "store.api/dto"

type CollectionService interface {
	GetAll(uint) []*dto.GetCollection
	Create(*dto.CreateCollection, uint) (*dto.GetCollection, error)
	AddCard(*dto.CreateCardSlot, uint, uint) (*dto.GetCollection, error)
}
