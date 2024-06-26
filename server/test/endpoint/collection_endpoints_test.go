package endpoint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
	"store.api/model"
)

func Test_Collection_ShouldCreate(t *testing.T) {
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

	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}

	// act
	w, body := req(r, t, "POST", "/api/v1/collection", data, token)

	var result dto.GetCollection
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 201, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, data.Name, result.Name)
	assert.Equal(t, data.Description, result.Description)
	assert.Len(t, result.Cards, 0)
}

func Test_Collection_ShouldNotCreateNotVerified(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")

	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/collection", data, token)

	// assert
	assert.Equal(t, 403, w.Code)
}

func Test_Collection_ShouldNotCreateEmptyName(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")

	data := dto.PostCollection{
		Name:        "",
		Description: "collection description",
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/collection", data, token)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Collection_ShouldNotCreateUnauthorized(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)

	data := dto.PostCollection{
		Name:        "collection",
		Description: "collection description",
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/collection", data, "")

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_Collection_ShouldFetchAll(t *testing.T) {
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

	req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection2",
		Description: "collection description",
	}, token)
	req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection3",
		Description: "collection description",
	}, token)

	// act
	w, body := req(r, t, "GET", "/api/v1/collection/all", nil, token)

	var result []*dto.GetCollection
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Len(t, result, 3)
}

func Test_Collection_ShouldNotFetchAllUnauthorized(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)

	// act
	w, _ := req(r, t, "GET", "/api/v1/collection/all", nil, "")

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_Collection_ShouldFetchById(t *testing.T) {
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

	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}

	_, colBody := req(r, t, "POST", "/api/v1/collection", data, token)

	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	// act
	w, body := req(r, t, "GET", fmt.Sprintf("/api/v1/collection/%v", collection.ID), nil, token)
	var result dto.GetCollection
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, collection.Name, result.Name)
	assert.Equal(t, collection.Description, result.Description)
	assert.Equal(t, collection.ID, result.ID)
	assert.Equal(t, collection.OwnerId, result.OwnerId)
	assert.ElementsMatch(t, collection.Cards, result.Cards)
}

func Test_Collection_ShouldNotFetchByIdNotFound(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")

	// act
	w, _ := req(r, t, "GET", "/api/v1/collection/12", nil, token)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Collection_ShouldNotFetchByIdUnauthorized(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)

	// act
	w, _ := req(r, t, "GET", "/api/v1/collection/12", nil, "")

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_Collection_ShouldNotAddCardUnverified(t *testing.T) {
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	data := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: 3,
	}
	err = db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", false).
		Error

	if err != nil {
		t.Fatal(err)
	}

	// act
	w, _ := req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	// assert
	assert.Equal(t, 403, w.Code)
}

func Test_Collection_ShouldAddCard(t *testing.T) {
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	data := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: 3,
	}

	// act
	w, body := req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	var result dto.GetCollection
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, collection.ID, result.ID)
	assert.Equal(t, collection.Name, result.Name)
	assert.Equal(t, collection.Description, result.Description)
	assert.Len(t, result.Cards, 1)
	assert.Equal(t, cardId, result.Cards[0].CardId)
}

func Test_Collection_ShouldNotEditSlotNegativeAmount(t *testing.T) {
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	data := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: -3,
	}

	// act
	w, _ := req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Collection_ShouldNotAddCardInvalidCollectionId(t *testing.T) {
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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

	data := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: 3,
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/collection/12", data, token)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Collection_ShouldNotAddCardInvalidCardId(t *testing.T) {
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	data := dto.PostCollectionSlot{
		CardId: 12,
		Amount: 3,
	}

	// act
	w, _ := req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Collection_ShouldNotAddCardUnauthorized(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)

	data := dto.PostCollectionSlot{
		CardId: 1,
		Amount: 3,
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/collection/1", data, "")

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_Collection_ShouldAddCardConsecutive(t *testing.T) {
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	data := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: 3,
	}

	req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	// act
	w, body := req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	var result dto.GetCollection
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, collection.ID, result.ID)
	assert.Equal(t, collection.Name, result.Name)
	assert.Equal(t, collection.Description, result.Description)
	assert.Len(t, result.Cards, 1)
	assert.Equal(t, cardId, result.Cards[0].CardId)
	assert.Equal(t, uint(data.Amount*2), result.Cards[0].Amount)
}

func Test_Collection_ShouldSubtractCardConsecutive(t *testing.T) {
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	data1 := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: 10,
	}
	data2 := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: -3,
	}

	req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data1, token)

	// act
	w, body := req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data2, token)

	var result dto.GetCollection
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, collection.ID, result.ID)
	assert.Equal(t, collection.Name, result.Name)
	assert.Equal(t, collection.Description, result.Description)
	assert.Len(t, result.Cards, 1)
	assert.Equal(t, cardId, result.Cards[0].CardId)
	assert.Equal(t, uint(data1.Amount+data2.Amount), result.Cards[0].Amount)
}

func Test_Collection_ShouldRemoveCard(t *testing.T) {
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	data := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: 3,
	}

	req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	data.Amount = -data.Amount
	// act
	w, body := req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	var result dto.GetCollection
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, collection.ID, result.ID)
	assert.Equal(t, collection.Name, result.Name)
	assert.Equal(t, collection.Description, result.Description)
	assert.Len(t, result.Cards, 0)
}

func Test_Collection_ShouldDeleteCollection(t *testing.T) {
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	// act
	deleteW, _ := req(r, t, "DELETE", fmt.Sprintf("/api/v1/collection/%d", collection.ID), nil, token)
	getW, _ := req(r, t, "GET", fmt.Sprintf("/api/v1/collection/%d", collection.ID), nil, token)

	// assert
	assert.Equal(t, 200, deleteW.Code)
	assert.Equal(t, 404, getW.Code)
}

func Test_Collection_ShouldUpdate(t *testing.T) {
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

	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	_, collectionData := req(r, t, "POST", "/api/v1/collection", data, token)

	var collection dto.GetCollection
	err = json.Unmarshal(collectionData, &collection)
	if err != nil {
		panic(err)
	}
	newData := dto.PostCollection{
		Name:        "collection2",
		Description: "collection description 1",
	}

	// act
	w, body := req(r, t, "PATCH", fmt.Sprintf("/api/v1/collection/%v", collection.ID), newData, token)
	var result dto.GetCollection
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, newData.Name, result.Name)
	assert.Equal(t, newData.Description, result.Description)
	assert.Len(t, result.Cards, 0)
}

func Test_Collection_ShouldNotUpdateBadData(t *testing.T) {
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

	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	_, collectionData := req(r, t, "POST", "/api/v1/collection", data, token)

	var collection dto.GetCollection
	err = json.Unmarshal(collectionData, &collection)
	if err != nil {
		panic(err)
	}
	newData := dto.PostCollection{
		Name:        "",
		Description: "collection description 1",
	}

	// act
	w, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/collection/%v", collection.ID), newData, token)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Collection_ShouldNotUpdateCollectionNotFound(t *testing.T) {
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

	newData := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description 1",
	}

	// act
	w, _ := req(r, t, "PATCH", "/api/v1/collection/1", newData, token)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Collection_ShouldNotUpdateUnverified(t *testing.T) {
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

	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	_, collectionData := req(r, t, "POST", "/api/v1/collection", data, token)

	var collection dto.GetCollection
	err = json.Unmarshal(collectionData, &collection)
	if err != nil {
		panic(err)
	}

	err = db.
		Model(&model.User{}).
		Where("username=?", username).
		Update("verified", false).
		Error

	if err != nil {
		t.Fatal(err)
	}

	newData := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description 1",
	}

	// act
	w, _ := req(r, t, "PATCH", fmt.Sprintf("/api/v1/collection/%v", collection.ID), newData, token)

	// assert
	assert.Equal(t, 403, w.Code)
}

// *** weird edge cases ***

func Test_Collection_ShouldUpdateCardSlotAfterCardModification(t *testing.T) {
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
		Create(&model.CardType{
			ID:       "CT1",
			LongName: "Card type 1",
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

	_, colBody := req(r, t, "POST", "/api/v1/collection", dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, token)
	var collection dto.GetCollection
	err = json.Unmarshal(colBody, &collection)
	if err != nil {
		t.Fatal(err)
	}

	data := dto.PostCollectionSlot{
		CardId: cardId,
		Amount: 3,
	}

	req(r, t, "POST", fmt.Sprintf("/api/v1/collection/%d", collection.ID), data, token)

	newCardName := "card2"
	err = db.
		Model(&model.Card{}).
		Where("name=?", "card1").
		UpdateColumn("name", newCardName).
		Error
	if err != nil {
		panic(err)
	}

	// act
	w, body := req(r, t, "GET", fmt.Sprintf("/api/v1/collection/%d", collection.ID), nil, token)

	var result dto.GetCollection
	err = json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, collection.ID, result.ID)
	assert.Len(t, result.Cards, 1)
	assert.Equal(t, cardId, result.Cards[0].CardId)
}
