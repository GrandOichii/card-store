package repository

import "store.api/model"

type CollectionRepository interface {
	FindByOwnerId(ownerId uint) []*model.Collection
	Save(*model.Collection) error
	FindById(id uint) *model.Collection
	Update(*model.Collection) error
	UpdateSlot(slot *model.CollectionSlot) error
	DeleteSlot(slot *model.CollectionSlot) error
	Delete(id uint) error
}
