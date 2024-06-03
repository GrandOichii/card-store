package cache

import (
	"store.api/model"
)

type ExpansionCache interface {
	Remember([]*model.Expansion)
	Forget()
	Get() []*model.Expansion
}

type NoExpansionCache struct {
}

func (c *NoExpansionCache) Remember([]*model.Expansion) {
}

func (c *NoExpansionCache) Forget() {
}

func (c *NoExpansionCache) Get() []*model.Expansion {
	return nil
}
