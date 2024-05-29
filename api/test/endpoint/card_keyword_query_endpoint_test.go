package endpoint_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"store.api/dto"
	"store.api/model"
)

func setupDb(t *testing.T, db *gorm.DB) {
	var err error
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
		Create(&model.CardType{
			ID:       "CT2",
			LongName: "Card type 2",
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
		Create(&model.CardKey{
			ID:      "key2",
			EngName: "card2",
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
}

func Test_CardKeywordQuery_ShouldFetchByName(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	setupDb(t, db)

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

	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card2",
		Text:      "card text",
		Price:     400,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key2",
		Expansion: "exp1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card2",
		Text:      "card text",
		Price:     400,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?t=ard1", nil, "")
	var cards []*dto.GetCard
	err = json.Unmarshal(body, &cards)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, cards, 2)
}

func Test_CardKeywordQuery_ShouldFetchByType(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	setupDb(t, db)

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

	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT2",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card2",
		Text:      "card text",
		Price:     400,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key2",
		Expansion: "exp1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card2",
		Text:      "card text",
		Price:     400,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?t=ct1", nil, "")
	var cards []*dto.GetCard
	err = json.Unmarshal(body, &cards)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, cards, 2)
}
