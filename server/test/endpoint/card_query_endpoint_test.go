package endpoint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
	"store.api/model"
	"store.api/service"
)

func Test_Card_ShouldFetchById(t *testing.T) {
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

	_, b := req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card1",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
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
	r, _ := setupRouter(10)

	// act
	w, _ := req(r, t, "GET", "/api/v1/card/1", nil, "")

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Card_ShouldFetchByType(t *testing.T) {
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
		Price:     10,
		Type:      "CT2",
		Language:  "ENG",
		Key:       "key2",
		Expansion: "exp1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?type=CT1", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 1)
	assert.Equal(t, "card1", queryResult.Cards[0].Name)
}

func Test_Card_ShouldFetchByName(t *testing.T) {
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
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key2",
		Expansion: "exp1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?name=d2", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 1)
	assert.Equal(t, "card2", queryResult.Cards[0].Name)
}

func Test_Card_ShouldFetchByMinPrice(t *testing.T) {
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
		Key:       "key1",
		Expansion: "exp1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?minPrice=300", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 1)
	assert.Equal(t, "card2", queryResult.Cards[0].Name)
}

func Test_Card_ShouldFetchByMaxPrice(t *testing.T) {
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
		Key:       "key1",
		Expansion: "exp1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?maxPrice=300", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 1)
	assert.Equal(t, "card1", queryResult.Cards[0].Name)
}

func Test_Card_ShouldFetchByLanguage(t *testing.T) {
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
		Create(&model.Language{
			ID:       "RUS",
			LongName: "Russian",
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
		Language:  "RUS",
		Key:       "key1",
		Expansion: "exp1",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?lang=ENG", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 1)
	assert.Equal(t, "card1", queryResult.Cards[0].Name)
}

func Test_Card_ShouldFetchByCardKey(t *testing.T) {
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
	w, body := req(r, t, "GET", "/api/v1/card?key=key1", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 2)
	assert.Equal(t, "key1", queryResult.Cards[0].Key)
	assert.Equal(t, "key1", queryResult.Cards[1].Key)
}

func Test_Card_ShouldFetchByExpansion(t *testing.T) {
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
		Key:       "key1",
		Expansion: "exp2",
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
	w, body := req(r, t, "GET", "/api/v1/card?expansion=exp1", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 2)
	assert.Equal(t, "exp1", queryResult.Cards[0].Expansion)
	assert.Equal(t, "expansion1", queryResult.Cards[0].ExpansionName)
	assert.Equal(t, "exp1", queryResult.Cards[1].Expansion)
	assert.Equal(t, "expansion1", queryResult.Cards[1].ExpansionName)
}

func Test_Card_ShouldFetchPages(t *testing.T) {
	// arrange
	r, db := setupRouter(2)

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
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card3",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, token)

	// act
	w1, body1 := req(r, t, "GET", "/api/v1/card?page=1", nil, "")
	var query1 service.CardQueryResult
	err1 := json.Unmarshal(body1, &query1)

	w2, body2 := req(r, t, "GET", "/api/v1/card?page=2", nil, "")
	var query2 service.CardQueryResult
	err2 := json.Unmarshal(body2, &query2)

	// assert
	assert.Equal(t, 200, w1.Code)
	assert.Nil(t, err1)
	assert.Len(t, query1.Cards, 2)

	assert.Equal(t, 200, w2.Code)
	assert.Nil(t, err2)
	assert.Len(t, query2.Cards, 1)
}

func Test_Card_ShouldFetchByBeingInStock(t *testing.T) {
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
			FullName:  "expansion1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:          "card1",
		Text:          "card text",
		Price:         10,
		Type:          "CT1",
		Language:      "ENG",
		Key:           "key1",
		Expansion:     "exp1",
		InStockAmount: 20,
	}, token)
	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:          "card2",
		Text:          "card text",
		Price:         400,
		Type:          "CT1",
		Language:      "ENG",
		Key:           "key1",
		Expansion:     "exp1",
		InStockAmount: 10,
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
	w, body := req(r, t, "GET", "/api/v1/card?inStockOnly=true", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 2)
}

func Test_Card_ShouldFetchFoilOnly(t *testing.T) {
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
			FullName:  "expansion1",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Foiling{
			ID:              "foil1",
			Label:           "Foil1",
			DescriptiveName: "Foiling",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}
	err = db.
		Create(&model.Foiling{
			ID:              "foil2",
			Label:           "Foil2",
			DescriptiveName: "Foiling",
		}).
		Error
	if err != nil {
		t.Fatal(err)
	}

	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:          "card1",
		Text:          "card text",
		Price:         10,
		Type:          "CT1",
		Language:      "ENG",
		Key:           "key1",
		Expansion:     "exp1",
		InStockAmount: 20,
	}, token)

	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:          "card2",
		Text:          "card text",
		Price:         400,
		Type:          "CT1",
		Language:      "ENG",
		Key:           "key1",
		Expansion:     "exp1",
		InStockAmount: 10,
		Foiling:       "foil1",
	}, token)

	req(r, t, "POST", "/api/v1/card", dto.PostCard{
		Name:      "card2",
		Text:      "card text",
		Price:     400,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
		Foiling:   "foil2",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/card?foilOnly=true", nil, "")
	var queryResult service.CardQueryResult
	err = json.Unmarshal(body, &queryResult)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, queryResult.Cards, 2)
}
