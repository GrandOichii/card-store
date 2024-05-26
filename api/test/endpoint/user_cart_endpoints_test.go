package endpoint_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
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
