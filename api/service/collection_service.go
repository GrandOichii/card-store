package service

import (
	"errors"

	"store.api/dto"
)

var ErrNotVerified = errors.New("user is not verified")

type CollectionService interface {
	GetAll(uint) []*dto.GetCollection
	Create(*dto.CreateCollection, uint) (*dto.GetCollection, error)
	EditCard(*dto.PostCardSlot, uint, uint) (*dto.GetCollection, error)
	GetById(uint, uint) (*dto.GetCollection, error)
}
