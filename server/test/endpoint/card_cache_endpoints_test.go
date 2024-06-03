package endpoint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"store.api/dto"
	"store.api/model"
	"store.api/service"
)

// * special tests for checking whether the cached information is updated when patching card price/availability/info

func create(t *testing.T, db *gorm.DB, data interface{}) {
	err := db.
		Create(data).Error
	if err != nil {
		t.Fatal(err)
	}
}

func Test_CardCache_ShouldCreateAndCache(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	create(t, db, &model.CardType{
		ID:        "ct1",
		LongName:  "CardType1",
		ShortName: "CT1",
	})
	create(t, db, &model.Language{
		ID:       "ENG",
		LongName: "English",
	})
	create(t, db, &model.CardKey{
		ID:      "key1",
		EngName: "card1",
	})
	create(t, db, &model.Expansion{
		ID:        "exp1",
		ShortName: "exp1",
		FullName:  "expansion",
	})

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "ct1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, cardBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(cardBody, &created)
	if err != nil {
		t.Fatal(err)
	}

	c := model.Card{}
	c.ID = created.ID
	err = db.
		Delete(&c).
		Error
	if err != nil {
		t.Fatal(err)
	}

	// act
	w, body := req(r, t, "GET", fmt.Sprintf("/api/v1/card/%v", created.ID), card, token)
	var result dto.GetCard
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, card.Name, result.Name)
	assert.Equal(t, card.Text, result.Text)
	assert.Equal(t, card.Price, result.Price)
	assert.Equal(t, card.Type, result.Type.ID)
	assert.Equal(t, card.Language, result.Language.ID)
}

func Test_CardCache_ShouldFetchWithPatchedPrice(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	create(t, db, &model.CardType{
		ID:        "ct1",
		LongName:  "CardType1",
		ShortName: "CT1",
	})
	create(t, db, &model.Language{
		ID:       "ENG",
		LongName: "English",
	})
	create(t, db, &model.CardKey{
		ID:      "key1",
		EngName: "card1",
	})
	create(t, db, &model.Expansion{
		ID:        "exp1",
		ShortName: "exp1",
		FullName:  "expansion",
	})

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "ct1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	create, cardBody := req(r, t, "POST", "/api/v1/card", card, token)
	assert.Equal(t, 201, create.Code)
	var created dto.GetCard
	err = json.Unmarshal(cardBody, &created)
	if err != nil {
		t.Fatal(err)
	}

	patch, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/price/%v", created.ID), dto.PriceUpdate{NewPrice: 100}, token)
	assert.Equal(t, 200, patch.Code)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?minPrice=9", nil, token)
	var result service.CardQueryResult
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 1)
}

func Test_CardCache_ShouldNotFetchWithPatchedPrice(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	create(t, db, &model.CardType{
		ID:        "ct1",
		LongName:  "CardType1",
		ShortName: "CT1",
	})
	create(t, db, &model.Language{
		ID:       "ENG",
		LongName: "English",
	})
	create(t, db, &model.CardKey{
		ID:      "key1",
		EngName: "card1",
	})
	create(t, db, &model.Expansion{
		ID:        "exp1",
		ShortName: "exp1",
		FullName:  "expansion",
	})

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "ct1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	create, cardBody := req(r, t, "POST", "/api/v1/card", card, token)
	assert.Equal(t, 201, create.Code)
	var created dto.GetCard
	err = json.Unmarshal(cardBody, &created)
	if err != nil {
		t.Fatal(err)
	}

	patch, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/price/%v", created.ID), dto.PriceUpdate{NewPrice: 100}, token)
	assert.Equal(t, 200, patch.Code)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?maxPrice=9", nil, token)
	var result service.CardQueryResult
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 0)
}

func Test_CardCache_ShouldFetchWithPatchedStockAmount(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	create(t, db, &model.CardType{
		ID:        "ct1",
		LongName:  "CardType1",
		ShortName: "CT1",
	})
	create(t, db, &model.Language{
		ID:       "ENG",
		LongName: "English",
	})
	create(t, db, &model.CardKey{
		ID:      "key1",
		EngName: "card1",
	})
	create(t, db, &model.Expansion{
		ID:        "exp1",
		ShortName: "exp1",
		FullName:  "expansion",
	})

	card := dto.PostCard{
		Name:          "card1",
		Text:          "card text",
		Price:         10,
		InStockAmount: 0,
		Type:          "ct1",
		Language:      "ENG",
		Key:           "key1",
		Expansion:     "exp1",
	}

	create, cardBody := req(r, t, "POST", "/api/v1/card", card, token)
	assert.Equal(t, 201, create.Code)
	var created dto.GetCard
	err = json.Unmarshal(cardBody, &created)
	if err != nil {
		t.Fatal(err)
	}

	patch, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/stocked/%v", created.ID), dto.StockedAmountUpdate{NewAmount: 10}, token)
	assert.Equal(t, 200, patch.Code)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?inStockOnly=true", nil, token)
	var result service.CardQueryResult
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 1)
}

func Test_CardCache_ShouldNotFetchWithPatchedStockAmount(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	create(t, db, &model.CardType{
		ID:        "ct1",
		LongName:  "CardType1",
		ShortName: "CT1",
	})
	create(t, db, &model.Language{
		ID:       "ENG",
		LongName: "English",
	})
	create(t, db, &model.CardKey{
		ID:      "key1",
		EngName: "card1",
	})
	create(t, db, &model.Expansion{
		ID:        "exp1",
		ShortName: "exp1",
		FullName:  "expansion",
	})

	card := dto.PostCard{
		Name:          "card1",
		Text:          "card text",
		Price:         10,
		InStockAmount: 10,
		Type:          "ct1",
		Language:      "ENG",
		Key:           "key1",
		Expansion:     "exp1",
	}

	create, cardBody := req(r, t, "POST", "/api/v1/card", card, token)
	assert.Equal(t, 201, create.Code)
	var created dto.GetCard
	err = json.Unmarshal(cardBody, &created)
	if err != nil {
		t.Fatal(err)
	}

	patch, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/stocked/%v", created.ID), dto.StockedAmountUpdate{NewAmount: 0}, token)
	assert.Equal(t, 200, patch.Code)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?inStockOnly=true", nil, token)
	var result service.CardQueryResult
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 0)
}

func Test_CardCache_ShouldFetchAfterUpdate(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	create(t, db, &model.CardType{
		ID:        "ct1",
		LongName:  "CardType1",
		ShortName: "CT1",
	})
	create(t, db, &model.Language{
		ID:       "ENG",
		LongName: "English",
	})
	create(t, db, &model.CardKey{
		ID:      "key1",
		EngName: "card1",
	})
	create(t, db, &model.Expansion{
		ID:        "exp1",
		ShortName: "exp1",
		FullName:  "expansion",
	})

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "ct1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, cardBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(cardBody, &created)
	if err != nil {
		t.Fatal(err)
	}

	update := card
	update.Name = "card2"

	req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/%v", created.ID), update, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?name=card2", update, token)
	var result service.CardQueryResult
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 1)
}

func Test_CardCache_ShouldNotFetchAfterUpdate(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	create(t, db, &model.CardType{
		ID:        "ct1",
		LongName:  "CardType1",
		ShortName: "CT1",
	})
	create(t, db, &model.Language{
		ID:       "ENG",
		LongName: "English",
	})
	create(t, db, &model.CardKey{
		ID:      "key1",
		EngName: "card1",
	})
	create(t, db, &model.Expansion{
		ID:        "exp1",
		ShortName: "exp1",
		FullName:  "expansion",
	})

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "ct1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, cardBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(cardBody, &created)
	if err != nil {
		t.Fatal(err)
	}

	update := card
	update.Name = "card2"

	req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/%v", created.ID), update, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?name=card1", update, token)
	var result service.CardQueryResult
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 0)
}

func Test_CardCache_ShouldFetchAndChangeAfterNewCard(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	create(t, db, &model.CardType{
		ID:        "ct1",
		LongName:  "CardType1",
		ShortName: "CT1",
	})
	create(t, db, &model.Language{
		ID:       "ENG",
		LongName: "English",
	})
	create(t, db, &model.CardKey{
		ID:      "key1",
		EngName: "card1",
	})
	create(t, db, &model.Expansion{
		ID:        "exp1",
		ShortName: "exp1",
		FullName:  "expansion",
	})

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "ct1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	query := "/api/v1/card?name=card"

	req(r, t, "POST", "/api/v1/card", card, token)

	req(r, t, "GET", query, nil, token)

	card.Name = "card2"
	req(r, t, "POST", "/api/v1/card", card, token)

	// act
	w, body := req(r, t, "GET", query, nil, token)
	var result service.CardQueryResult
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 2)
}
