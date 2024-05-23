package repository

import "store.api/model"

type CollectionRepository interface {
	FindByOwnerId(ownerId uint) []*model.Collection
}
