package cache

import (
	"context"
	"encoding/json"

	"github.com/valkey-io/valkey-go"
	"store.api/model"
)

const languageCacheKey = "languages"

type LanguageValkeyCache struct {
	client valkey.Client
}

func NewLanguageValkeyCache(client valkey.Client) *LanguageValkeyCache {
	return &LanguageValkeyCache{
		client: client,
	}
}

func (c *LanguageValkeyCache) Remember(languages []*model.Language) {
	json, err := json.Marshal(languages)
	if err != nil {
		panic(err)
	}
	err = c.client.Do(context.Background(), c.client.
		B().
		Set().
		Key(languageCacheKey).
		Value(string(json)).
		Build()).
		Error()
	if err != nil {
		panic(err)
	}
}

func (c *LanguageValkeyCache) Forget() {
	err := c.client.Do(context.Background(), c.client.
		B().
		Del().
		Key(languageCacheKey).
		Build()).Error()
	if err != nil {
		panic(err)
	}
}

func (c *LanguageValkeyCache) Get() []*model.Language {
	get := c.client.Do(context.Background(), c.client.
		B().
		Get().
		Key(languageCacheKey).
		Build())
	err := get.Error()
	if err != nil {
		if err == valkey.Nil {
			return nil
		}
		panic(err)
	}
	var result []*model.Language
	err = get.DecodeJSON(&result)
	if err != nil {
		panic(err)
	}
	return result
}
