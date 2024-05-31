package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valkey-io/valkey-go"
	"store.api/model"
)

type CardValkeyCache struct {
	client valkey.Client
}

func NewCardValkeyCache(client valkey.Client) *CardValkeyCache {
	return &CardValkeyCache{
		client: client,
	}
}

func (c *CardValkeyCache) ToKey(id uint) string {
	return fmt.Sprintf("card-%v", id)
}

func (c *CardValkeyCache) Remember(card *model.Card) {
	json, err := json.Marshal(card)
	if err != nil {
		panic(err)
	}
	err = c.client.Do(context.Background(), c.client.
		B().
		Set().
		Key(c.ToKey(card.ID)).
		Value(string(json)).
		Build()).
		Error()
	if err != nil {
		panic(err)
	}
}

func (c *CardValkeyCache) Forget(id uint) {
	err := c.client.Do(context.Background(), c.client.
		B().
		Del().
		Key(c.ToKey(id)).
		Build()).Error()
	if err != nil {
		panic(err)
	}
}

func (c *CardValkeyCache) Get(id uint) *model.Card {
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
	var result model.Card
	err = get.DecodeJSON(&result)
	if err != nil {
		panic(err)
	}
	return &result
}
