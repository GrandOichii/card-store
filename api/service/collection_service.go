package service

import "store.api/dto"

type CollectionService interface {
	GetAll(userId uint) []*dto.GetCollection
}
