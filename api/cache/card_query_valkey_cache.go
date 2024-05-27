package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valkey-io/valkey-go"
	"store.api/model"
)

type CardQueryValkeyCache struct {
	client valkey.Client
}

func NewCardQueryValkeyCache(client valkey.Client) *CardQueryValkeyCache {
	return &CardQueryValkeyCache{
		client: client,
	}
}

func (c *CardQueryValkeyCache) ToKey(rawQuery string) string {
	return fmt.Sprintf("cardQuery-%v", rawQuery)
}

func (c *CardQueryValkeyCache) Remember(rawQuery string, cardQuery []*model.Card) {
	json, err := json.Marshal(cardQuery)
	if err != nil {
		panic(err)
	}
	err = c.client.Do(context.Background(), c.client.
		B().
		Set().
		Key(c.ToKey(rawQuery)).
		Value(string(json)).
		Build()).
		Error()
	if err != nil {
		panic(err)
	}
}

func (c *CardQueryValkeyCache) Forget(rawQuery string) {
	err := c.client.Do(context.Background(), c.client.
		B().
		Del().
		Key(c.ToKey(rawQuery)).
		Build()).Error()
	if err != nil {
		panic(err)
	}
}

func (c *CardQueryValkeyCache) ForgetAll() {
	err := c.client.Do(context.Background(), c.client.
		B().
		Del().
		Key(c.ToKey("*")).
		Build()).Error()
	if err != nil {
		panic(err)
	}
}

func (c *CardQueryValkeyCache) Get(rawQuery string) []*model.Card {
	get := c.client.Do(context.Background(), c.client.
		B().
		Get().
		Key(c.ToKey(rawQuery)).
		Build())
	err := get.Error()
	if err != nil {
		if err == valkey.Nil {
			return nil
		}
		panic(err)
	}
	var result []*model.Card
	err = get.DecodeJSON(&result)
	if err != nil {
		panic(err)
	}
	return result
}
