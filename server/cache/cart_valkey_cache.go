package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valkey-io/valkey-go"
	"store.api/model"
)

type CartValkeyCache struct {
	client valkey.Client
}

func NewCartValkeyCache(client valkey.Client) *CartValkeyCache {
	return &CartValkeyCache{
		client: client,
	}
}

func (c *CartValkeyCache) ToKey(userId uint) string {
	return fmt.Sprintf("cart-%v", userId)
}

func (c *CartValkeyCache) Remember(cart *model.Cart) {
	json, err := json.Marshal(cart)
	if err != nil {
		panic(err)
	}
	err = c.client.Do(context.Background(), c.client.
		B().
		Set().
		Key(c.ToKey(cart.UserID)).
		Value(string(json)).
		Build()).
		Error()
	if err != nil {
		panic(err)
	}
}

func (c *CartValkeyCache) Forget(userId uint) {
	err := c.client.Do(context.Background(), c.client.
		B().
		Del().
		Key(c.ToKey(userId)).
		Build()).Error()
	if err != nil {
		panic(err)
	}
}

func (c *CartValkeyCache) Get(userId uint) *model.Cart {
	get := c.client.Do(context.Background(), c.client.
		B().
		Get().
		Key(c.ToKey(userId)).
		Build())
	err := get.Error()
	if err != nil {
		if err == valkey.Nil {
			return nil
		}
		panic(err)
	}
	var result model.Cart
	err = get.DecodeJSON(&result)
	if err != nil {
		panic(err)
	}
	return &result
}
