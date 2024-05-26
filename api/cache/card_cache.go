package cache

import (
	"store.api/model"
)

type CardCache interface {
	Remember(*model.Card)
	Forget(uint)
	Get(uint) *model.Card
}

type NoCardCache struct {
}

func (c *NoCardCache) Remember(*model.Card) {
}

func (c *NoCardCache) Forget(uint) {
}

func (c *NoCardCache) Get(uint) *model.Card {
	return nil
}
