package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valkey-io/valkey-go"
	"store.api/model"
)

type CollectionValkeyCache struct {
	client valkey.Client
}

func NewCollectionValkeyCache(client valkey.Client) *CollectionValkeyCache {
	return &CollectionValkeyCache{
		client: client,
	}
}

func (c *CollectionValkeyCache) ToKey(id uint) string {
	return fmt.Sprintf("collection-%v", id)
}

func (c *CollectionValkeyCache) Remember(Collection *model.Collection) {
	json, err := json.Marshal(Collection)
	if err != nil {
		panic(err)
	}
	err = c.client.Do(context.Background(), c.client.
		B().
		Set().
		Key(c.ToKey(Collection.ID)).
		Value(string(json)).
		Build()).
		Error()
	if err != nil {
		panic(err)
	}
}

func (c *CollectionValkeyCache) Forget(id uint) {
	err := c.client.Do(context.Background(), c.client.
		B().
		Del().
		Key(c.ToKey(id)).
		Build()).Error()
	if err != nil {
		panic(err)
	}
}

func (c *CollectionValkeyCache) Get(id uint) *model.Collection {
	get := c.client.Do(context.Background(), c.client.
		B().
		Get().
		Key(c.ToKey(id)).
		Build())
	err := get.Error()
	if err != nil {
		if err == valkey.Nil {
			return nil
		}
		panic(err)
	}
	var result model.Collection
	err = get.DecodeJSON(&result)
	if err != nil {
		panic(err)
	}
	return &result
}
