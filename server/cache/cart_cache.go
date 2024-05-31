package cache

import (
	"store.api/model"
)

type CartCache interface {
	Remember(*model.Cart)
	Forget(uint)
	Get(uint) *model.Cart
}

type NoCartCache struct {
}

func (c *NoCartCache) Remember(*model.Cart) {
}

func (c *NoCartCache) Forget(uint) {
}

func (c *NoCartCache) Get(uint) *model.Cart {
	return nil
}
