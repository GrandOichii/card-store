package endpoint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
	"store.api/model"
)

// TODO? add more detailed checks

func Test_Card_ShouldNotCreate(t *testing.T) {
	// arrange
	r, _ := setupRouter()

	// act
	w, _ := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
	}, "")

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_Card_ShouldCreate(t *testing.T) {
	// arrange
	r, db := setupRouter()
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
		Model(&model.CardType{}).
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	// act
	w, body := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card1",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, token)
	var result dto.GetCard
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 201, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, "card1", result.Name)
}

func Test_Card_ShouldNotCreateNotEnoughPrivileges(t *testing.T) {
	// arrange
	r, db := setupRouter()
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.CardType{}).
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
			w, _ := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
				Name:  "card name",
				Text:  "card text",
				Price: 10,
				Type:  "CT1",
			}, token)

			// assert
			assert.Equal(t, 403, w.Code)
		})
	}
}

func Test_Card_ShouldFetchById(t *testing.T) {
	// arrange
	r, db := setupRouter()
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
		Model(&model.CardType{}).
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	_, b := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card1",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, token)

	var created dto.GetCard
	err = json.Unmarshal(b, &created)
	if err != nil {
		t.Fatal(err)
	}

	// act
	w, body := req(r, t, "GET", "/api/v1/card/"+fmt.Sprint(created.ID), nil, "")
	var result dto.GetCard
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, "card1", result.Name)
}

func Test_Card_ShouldNotFetchById(t *testing.T) {
	// arrange
	r, _ := setupRouter()

	// act
	w, _ := req(r, t, "GET", "/api/v1/card/1", nil, "")

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Card_ShouldFetchByType(t *testing.T) {
	// arrange
	r, db := setupRouter()

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
		Model(&model.CardType{}).
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Model(&model.CardType{}).
		Create(&model.CardType{
			ID:       "CT2",
			LongName: "Card type 2",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card1",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card2",
		Text:  "card text",
		Price: 10,
		Type:  "CT2",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?type=CT1", nil, "")
	var cards []*dto.GetCard
	err = json.Unmarshal(body, &cards)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, cards, 1)
	assert.Equal(t, "card1", cards[0].Name)
}

func Test_ShouldFetchByName(t *testing.T) {
	// arrange
	r, db := setupRouter()

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
		Model(&model.CardType{}).
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card1",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card2",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?name=d2", nil, "")
	var cards []*dto.GetCard
	err = json.Unmarshal(body, &cards)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, cards, 1)
	assert.Equal(t, "card2", cards[0].Name)
}

func Test_ShouldFetchByMinPrice(t *testing.T) {
	// arrange
	r, db := setupRouter()

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
		Model(&model.CardType{}).
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card1",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card2",
		Text:  "card text",
		Price: 400,
		Type:  "CT1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?minPrice=300", nil, "")
	var cards []*dto.GetCard
	err = json.Unmarshal(body, &cards)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, cards, 1)
	assert.Equal(t, "card2", cards[0].Name)
}

func Test_ShouldFetchByMaxPrice(t *testing.T) {
	// arrange
	r, db := setupRouter()

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
		Model(&model.CardType{}).
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card1",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card2",
		Text:  "card text",
		Price: 400,
		Type:  "CT1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?maxPrice=300", nil, "")
	var cards []*dto.GetCard
	err = json.Unmarshal(body, &cards)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, cards, 1)
	assert.Equal(t, "card1", cards[0].Name)
}
