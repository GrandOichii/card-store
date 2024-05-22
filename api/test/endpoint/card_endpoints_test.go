package endpoint_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
	"store.api/model"
)

func Test_ShouldFetchAll(t *testing.T) {
	// arrange
	r, _ := setupRouter()

	// act
	w, _ := req(r, t, "GET", "/api/v1/card/all", nil, "")

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_ShouldNotCreate(t *testing.T) {
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

func Test_ShouldCreate(t *testing.T) {
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
		panic(err)
	}

	// act
	w, _ := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
	}, token)

	// assert
	assert.Equal(t, 201, w.Code)
}

func Test_ShouldNotCreateNotEnoughPrivileges(t *testing.T) {
	// arrange
	r, db := setupRouter()
	username := "user"
	token := loginAs(r, t, username, "password", "mail@mail.com")
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
				panic(err)
			}

			// act
			w, _ := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
				Name:  "card name",
				Text:  "card text",
				Price: 10,
			}, token)

			// assert
			assert.Equal(t, 403, w.Code)
		})
	}
}

func Test_ShouldFetchById(t *testing.T) {
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
		panic(err)
	}

	_, b := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
	}, token)
	var created dto.GetCard
	err = json.Unmarshal(b, &created)
	if err != nil {
		panic(err)
	}

	// act
	w, _ := req(r, t, "GET", "/api/v1/card/"+fmt.Sprint(created.ID), nil, "")

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_ShouldNotFetchById(t *testing.T) {
	// arrange
	r, _ := setupRouter()

	// act
	w, _ := req(r, t, "GET", "/api/v1/card/1", nil, "")

	// assert
	assert.Equal(t, 404, w.Code)
}
