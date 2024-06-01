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