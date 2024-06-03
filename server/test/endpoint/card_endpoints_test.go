package endpoint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
	"store.api/model"
)

func Test_Card_ShouldNotCreateNoType(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)

	// act
	w, _ := req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
	}, "")

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_Card_ShouldCreate(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	// act
	w, body := req(r, t, "POST", "/api/v1/card", card, token)
	var result dto.GetCard
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 201, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, card.Name, result.Name)
	assert.Equal(t, card.Text, result.Text)
	assert.Equal(t, card.Price, result.Price)
	assert.Equal(t, card.Type, result.Type.ID)
	assert.Equal(t, card.Language, result.Language.ID)
}

func Test_Card_ShouldNotCreateNotEnoughPrivileges(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc       string
		isAdmin    bool
		isVerified bool
	}{
		{
			desc:       "Not admin, not verified",
			isAdmin:    false,
			isVerified: false,
		},
		{
			desc:       "Admin, not verified",
			isAdmin:    true,
			isVerified: false,
		},
		{
			desc:       "Not admin, verified",
			isAdmin:    false,
			isVerified: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := db.
				Model(&model.User{}).
				Where("username=?", username).
				Update("is_admin", tC.isAdmin).
				Update("verified", tC.isVerified).
				Error

			if err != nil {
				t.Fatal(err)
			}

			// act
			w, _ := req(r, t, "POST", "/api/v1/card", dto.PostCard{
				Name:      "card name",
				Text:      "card text",
				Price:     10,
				Type:      "CT1",
				Language:  "ENG",
				Key:       "key1",
				Expansion: "exp1",
			}, token)

			// assert
			assert.Equal(t, 403, w.Code)
		})
	}
}

func Test_Card_ShouldPatch(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	update := card
	update.Name = "card2"

	// act
	w, body := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/%v", created.ID), update, token)
	var result dto.GetCard
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, update.Name, result.Name)
	assert.Equal(t, update.Text, result.Text)
	assert.Equal(t, update.Price, result.Price)
	assert.Equal(t, update.Type, result.Type.ID)
	assert.Equal(t, update.Language, result.Language.ID)
}

func Test_Card_ShouldNotPatchBadData1(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Language:  "ENG",
		Type:      "CT1",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	update := card
	update.Name = ""

	// act
	w, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/%v", created.ID), update, token)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Card_ShouldNotPatchUnauthorized(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	update := card
	update.Name = "card2"

	err = db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", false).
		Error

	if err != nil {
		t.Fatal(err)
	}

	// act
	w, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/%v", created.ID), update, token)

	// assert
	assert.Equal(t, 403, w.Code)
}

func Test_Card_ShouldNotPatchCardNotFound(t *testing.T) {
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

	update := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	// act
	w, _ := req(r, t, "PATCH", "/api/v1/card/1", update, token)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Card_ShouldNotPatchBadData2(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	update := card
	update.Language = "RUS"

	// act
	w, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/%v", created.ID), update, token)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Card_ShouldPatchPrice(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	update := dto.PriceUpdate{
		NewPrice: 100,
	}

	// act
	w, body := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/price/%v", created.ID), update, token)
	var result dto.GetCard
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, card.Name, result.Name)
	assert.Equal(t, card.Text, result.Text)
	assert.Equal(t, update.NewPrice, result.Price)
	assert.Equal(t, card.Type, result.Type.ID)
	assert.Equal(t, card.Language, result.Language.ID)
}

func Test_Card_ShouldNotPatchPriceCardNotFound(t *testing.T) {
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

	update := dto.PriceUpdate{
		NewPrice: 100,
	}

	// act
	w, _ := req(r, t, "PATCH", "/api/v1/card/price/1", update, token)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Card_ShouldPatchPriceBadRequest(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	update := dto.PriceUpdate{
		NewPrice: -100,
	}

	// act
	w, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/price/%v", created.ID), update, token)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Card_ShouldPatchPriceNotAuthorized(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	err = db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", false).
		Error

	if err != nil {
		t.Fatal(err)
	}

	update := dto.PriceUpdate{
		NewPrice: 100,
	}

	// act
	w, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/price/%v", created.ID), update, token)

	// assert
	assert.Equal(t, 403, w.Code)
}

func Test_Card_ShouldPatchInStockAmount(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	update := dto.StockedAmountUpdate{
		NewAmount: 100,
	}

	// act
	w, body := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/stocked/%v", created.ID), update, token)
	var result dto.GetCard
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, card.Name, result.Name)
	assert.Equal(t, card.Text, result.Text)
	assert.Equal(t, update.NewAmount, result.InStockAmount)
	assert.Equal(t, card.Type, result.Type.ID)
	assert.Equal(t, card.Language, result.Language.ID)
}

func Test_Card_ShouldNotPatchInStockAmountCardNotFound(t *testing.T) {
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

	update := dto.StockedAmountUpdate{
		NewAmount: 100,
	}

	// act
	w, _ := req(r, t, "PATCH", "/api/v1/card/InStockAmount/1", update, token)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Card_ShouldPatchInStockAmountNotAuthorized(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}

	_, createdBody := req(r, t, "POST", "/api/v1/card", card, token)
	var created dto.GetCard
	err = json.Unmarshal(createdBody, &created)
	if err != nil {
		panic(err)
	}

	err = db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", false).
		Error

	if err != nil {
		t.Fatal(err)
	}

	update := dto.StockedAmountUpdate{
		NewAmount: 100,
	}

	// act
	w, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/card/stocked/%v", created.ID), update, token)

	// assert
	assert.Equal(t, 403, w.Code)
}

func Test_Card_ShouldFetchLanguages(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "RU",
			LongName: "Russian",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	// act
	w, body := req(r, t, "GET", "/api/v1/card/languages", nil, "")
	var result []model.Language
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "ENG", result[0].ID)
	assert.Equal(t, "English", result[0].LongName)
	assert.Equal(t, "RU", result[1].ID)
	assert.Equal(t, "Russian", result[1].LongName)
}

func Test_Card_ShouldFetchExpansions(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("is_admin", true).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp2",
			ShortName: "exp2",
			FullName:  "expansion2",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	// act
	w, body := req(r, t, "GET", "/api/v1/card/expansions", nil, "")
	var result []model.Expansion
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "exp1", result[0].ID)
	assert.Equal(t, "exp1", result[0].ShortName)
	assert.Equal(t, "expansion1", result[0].FullName)
	assert.Equal(t, "exp2", result[1].ID)
	assert.Equal(t, "exp2", result[1].ShortName)
	assert.Equal(t, "expansion2", result[1].FullName)
}

func Test_Card_ShouldCreateWithFoiling(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Foiling{
			ID:              "foil1",
			Label:           "Foil",
			DescriptiveName: "Standard foil",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
		Foiling:   "foil1",
	}

	// act
	w, body := req(r, t, "POST", "/api/v1/card", card, token)
	var result dto.GetCard
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 201, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, card.Name, result.Name)
	assert.Equal(t, card.Text, result.Text)
	assert.Equal(t, card.Price, result.Price)
	assert.Equal(t, card.Type, result.Type.ID)
	assert.Equal(t, card.Language, result.Language.ID)
}

func Test_Card_ShouldNotCreateInvalidFoiling(t *testing.T) {
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

	err = db.
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.CardKey{
			ID:      "key1",
			EngName: "card1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Expansion{
			ID:        "exp1",
			ShortName: "exp1",
			FullName:  "expansion",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	card := dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
		Foiling:   "foil1",
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/card", card, token)

	// assert
	assert.Equal(t, 400, w.Code)
}
