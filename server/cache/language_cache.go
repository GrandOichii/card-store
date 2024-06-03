package cache

import (
	"store.api/model"
)

type LanguageCache interface {
	Remember([]*model.Language)
	Forget()
	Get() []*model.Language
}

type NoLanguageCache struct {
}

func (c *NoLanguageCache) Remember([]*model.Language) {
}

func (c *NoLanguageCache) Forget() {
}

func (c *NoLanguageCache) Get() []*model.Language {
	return nil
}
