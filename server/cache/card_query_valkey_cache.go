package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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

func (c *CardQueryValkeyCache) ToCountKey(rawQuery string) string {
	return fmt.Sprintf("cardQueryTotalCount-%v", rawQuery)
}

func (c *CardQueryValkeyCache) Remember(rawQuery string, cardQuery []*model.Card, amount int64) {
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

	err = c.client.Do(context.Background(), c.client.
		B().
		Set().
		Key(c.ToCountKey(rawQuery)).
		Value(strconv.FormatInt(amount, 10)).
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

func (c *CardQueryValkeyCache) Get(rawQuery string) ([]*model.Card, int64) {
	get := c.client.Do(context.Background(), c.client.
		B().
		Get().
		Key(c.ToKey(rawQuery)).
		Build())
	err := get.Error()
	if err != nil {
		if err == valkey.Nil {
			return nil, -1
		}
		panic(err)
	}
	getCount := c.client.Do(context.Background(), c.client.
		B().
		Get().
		Key(c.ToCountKey(rawQuery)).
		Build())
	err = get.Error()
	if err != nil {
		if err == valkey.Nil {
			return nil, -1
		}
		panic(err)
	}

	var result []*model.Card
	err = get.DecodeJSON(&result)
	if err != nil {
		panic(err)
	}

	count, err := getCount.AsInt64()
	if err != nil {
		panic(err)
	}

	return result, count
}
