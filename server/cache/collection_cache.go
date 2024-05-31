package cache

import (
	"store.api/model"
)

type CollectionCache interface {
	Remember(*model.Collection)
	Forget(uint)
	Get(uint) *model.Collection
}

type NoCollectionCache struct {
}

func (c *NoCollectionCache) Remember(*model.Collection) {
}

func (c *NoCollectionCache) Forget(uint) {
}

func (c *NoCollectionCache) Get(uint) *model.Collection {
	return nil
}
