package cache

import (
	"store.api/model"
)

type CardQueryCache interface {
	Remember(string, []*model.Card)
	Forget(string)
	ForgetAll()
	Get(string) []*model.Card
}

type NoCardQueryCache struct {
}

func (c *NoCardQueryCache) Remember(string, []*model.Card) {
}

func (c *NoCardQueryCache) Forget(string) {
}

func (c *NoCardQueryCache) Get(string) []*model.Card {
	return make([]*model.Card, 0)
}
