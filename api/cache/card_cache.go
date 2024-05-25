package cache

import (
	"store.api/dto"
)

type CardCache interface {
	Remember(*dto.GetCard) error
	Forget(uint) error
	Get(uint) (*dto.GetCard, error)
}

type NoCardCache struct {
}

func (c *NoCardCache) Remember(*dto.GetCard) error {
	return nil
}

func (c *NoCardCache) Forget(uint) error {
	return nil
}

func (c *NoCardCache) Get(uint) (*dto.GetCard, error) {
	return nil, nil
}
