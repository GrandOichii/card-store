package service

import (
	"errors"

	"store.api/dto"
)

var (
	ErrNotVerified        = errors.New("user is not verified")
	ErrCollectionNotFound = errors.New("collection not found")
)

type CollectionService interface {
	GetAll(uint) []*dto.GetCollection
	Create(*dto.PostCollection, uint) (*dto.GetCollection, error)
	EditCard(*dto.PostCardSlot, uint, uint) (*dto.GetCollection, error)
	GetById(uint, uint) (*dto.GetCollection, error)
	Delete(uint, uint) error
	UpdateInfo(*dto.PostCollection, uint, uint) (*dto.GetCollection, error)
}
