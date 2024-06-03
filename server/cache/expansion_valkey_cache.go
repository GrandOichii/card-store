package cache

import (
	"context"
	"encoding/json"

	"github.com/valkey-io/valkey-go"
	"store.api/model"
)

const expansionCacheKey = "expansions"

type ExpansionValkeyCache struct {
	client valkey.Client
}

func NewExpansionValkeyCache(client valkey.Client) *ExpansionValkeyCache {
	return &ExpansionValkeyCache{
		client: client,
	}
}

func (c *ExpansionValkeyCache) Remember(expansions []*model.Expansion) {
	json, err := json.Marshal(expansions)
	if err != nil {
		panic(err)
	}
	err = c.client.Do(context.Background(), c.client.
		B().
		Set().
		Key(expansionCacheKey).
		Value(string(json)).
		Build()).
		Error()
	if err != nil {
		panic(err)
	}
}

func (c *ExpansionValkeyCache) Forget() {
	err := c.client.Do(context.Background(), c.client.
		B().
		Del().
		Key(expansionCacheKey).
		Build()).Error()
	if err != nil {
		panic(err)
	}
}

func (c *ExpansionValkeyCache) Get() []*model.Expansion {
	get := c.client.Do(context.Background(), c.client.
		B().
		Get().
		Key(expansionCacheKey).
		Build())
	err := get.Error()
	if err != nil {
		if err == valkey.Nil {
			return nil
		}
		panic(err)
	}
	var result []*model.Expansion
	err = get.DecodeJSON(&result)
	if err != nil {
		panic(err)
	}
	return result
}
