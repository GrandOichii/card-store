package cache

import (
	"store.api/model"
)

type CardQueryCache interface {
	Remember(string, []*model.Card, int64)
	Forget(string)
	ForgetAll()
	Get(string) ([]*model.Card, int64)
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
