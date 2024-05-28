package endpoint_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
	"store.api/model"
)

func Test_UserCart_ShouldFetchCart(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)
	token := loginAs(r, t, "user1", "password", "mail@mail.com")

	// act
	w, body := req(r, t, "GET", "/api/v1/user/cart", nil, token)
	var result dto.GetCart
	err := json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 0)
}

func Test_UserCart_ShouldNotFetchCart(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)
	createUser(r, t, "user1", "password", "mail@mail.com")

	// act
	w, _ := req(r, t, "GET", "/api/v1/user/cart", nil, "")

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_UserCart_ShouldAddCard(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", true).
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error

	if err != nil {
		t.Fatal(err)
	}
	adminId := createAdmin(r, t, db)
	cardId := createCard(t, db, &model.Card{
		Name:        "card1",
		Text:        "card text",
		Price:       1,
		PosterID:    adminId,
		CardTypeID:  "CT1",
		LanguageID:  "ENG",
		CardKeyID:   "key1",
		ExpansionID: "exp1",
	})

	data := dto.PostCartSlot{
		CardId: cardId,
		Amount: 3,
	}

	// act
	w, body := req(r, t, "POST", "/api/v1/user/cart", data, token)

	var result dto.GetCart
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 1)
	assert.Equal(t, cardId, result.Cards[0].CardId)
}

func Test_UserCart_ShouldNotEditSlotNegativeAmount(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", true).
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error

	if err != nil {
		t.Fatal(err)
	}

	adminId := createAdmin(r, t, db)
	cardId := createCard(t, db, &model.Card{
		Name:        "card1",
		Text:        "card text",
		Price:       1,
		PosterID:    adminId,
		CardTypeID:  "CT1",
		LanguageID:  "ENG",
		CardKeyID:   "key1",
		ExpansionID: "exp1",
	})

	data := dto.PostCartSlot{
		CardId: cardId,
		Amount: -3,
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/user/cart", data, token)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_UserCart_ShouldNotAddCardInvalidCardId(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", true).
		Error

	if err != nil {
		t.Fatal(err)
	}

	data := dto.PostCartSlot{
		CardId: 12,
		Amount: 3,
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/user/cart", data, token)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_UserCart_ShouldNotAddCardUnauthorized(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)

	data := dto.PostCartSlot{
		CardId: 1,
		Amount: 3,
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/user/cart", data, "")

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_UserCart_ShouldAddCardConsecutive(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", true).
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error

	if err != nil {
		t.Fatal(err)
	}
	adminId := createAdmin(r, t, db)
	cardId := createCard(t, db, &model.Card{
		Name:        "card1",
		Text:        "card text",
		Price:       1,
		PosterID:    adminId,
		CardTypeID:  "CT1",
		LanguageID:  "ENG",
		CardKeyID:   "key1",
		ExpansionID: "exp1",
	})

	data := dto.PostCartSlot{
		CardId: cardId,
		Amount: 3,
	}

	req(r, t, "POST", "/api/v1/user/cart", data, token)

	// act
	w, body := req(r, t, "POST", "/api/v1/user/cart", data, token)

	var result dto.GetCart
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 1)
	assert.Equal(t, cardId, result.Cards[0].CardId)
	assert.Equal(t, uint(data.Amount*2), result.Cards[0].Amount)
}

func Test_User_Cart_ShouldSubtractCardConsecutive(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", true).
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error

	if err != nil {
		t.Fatal(err)
	}
	adminId := createAdmin(r, t, db)
	cardId := createCard(t, db, &model.Card{
		Name:        "card1",
		Text:        "card text",
		Price:       1,
		PosterID:    adminId,
		CardTypeID:  "CT1",
		LanguageID:  "ENG",
		CardKeyID:   "key1",
		ExpansionID: "exp1",
	})

	data1 := dto.PostCartSlot{
		CardId: cardId,
		Amount: 10,
	}
	data2 := dto.PostCartSlot{
		CardId: cardId,
		Amount: -3,
	}

	req(r, t, "POST", "/api/v1/user/cart", data1, token)

	// act
	w, body := req(r, t, "POST", "/api/v1/user/cart", data2, token)

	var result dto.GetCart
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 1)
	assert.Equal(t, cardId, result.Cards[0].CardId)
	assert.Equal(t, uint(data1.Amount+data2.Amount), result.Cards[0].Amount)
}

func Test_UserCart_ShouldRemoveCard(t *testing.T) {
	// arrange
	r, db := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
	err := db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", true).
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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
		Create(&model.Language{
			ID:       "ENG",
			LongName: "English",
		}).
		Error

	if err != nil {
		t.Fatal(err)
	}
	adminId := createAdmin(r, t, db)
	cardId := createCard(t, db, &model.Card{
		Name:        "card1",
		Text:        "card text",
		Price:       1,
		PosterID:    adminId,
		CardTypeID:  "CT1",
		LanguageID:  "ENG",
		CardKeyID:   "key1",
		ExpansionID: "exp1",
	})

	data := dto.PostCartSlot{
		CardId: cardId,
		Amount: 3,
	}

	req(r, t, "POST", "/api/v1/user/cart", data, token)

	data.Amount = -data.Amount
	// act
	w, body := req(r, t, "POST", "/api/v1/user/cart", data, token)

	var result dto.GetCart
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result.Cards, 0)
}
