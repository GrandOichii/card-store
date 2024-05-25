package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valkey-io/valkey-go"
	"store.api/dto"
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

func (c *CardValkeyCache) Remember(card *dto.GetCard) error {
	json, err := json.Marshal(card)
	if err != nil {
		return err
	}
	return c.client.Do(context.Background(), c.client.
		B().
		Set().
		Key(c.ToKey(card.ID)).
		Value(string(json)).
		Build()).
		Error()

}

func (c *CardValkeyCache) Forget(id uint) error {
	return c.client.Do(context.Background(), c.client.
		B().
		JsonForget().
		Key(c.ToKey(id)).
		Build()).Error()
}

func (c *CardValkeyCache) Get(id uint) (*dto.GetCard, error) {
	get := c.client.Do(context.Background(), c.client.
		B().
		Get().
		Key(c.ToKey(id)).
		Build())
	err := get.Error()
	if err != nil {
		if err == valkey.Nil {
			return nil, nil
		}
		return nil, err
	}
	var result dto.GetCard
	err = get.DecodeJSON(&result)
	return &result, err
}
