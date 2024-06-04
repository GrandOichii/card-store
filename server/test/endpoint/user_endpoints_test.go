package endpoint_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
)

func Test_User_ShouldFetchInfo(t *testing.T) {
	// arrange
	r, _ := setupRouter(10)
	username := "user1"

	token := loginAs(r, t, username, "password", "mail@mail.com")

	// act
	w, body := req(r, t, "GET", "/api/v1/user", nil, token)
	var result dto.PrivateUserInfo
	err := json.Unmarshal(body, &result)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, err)
	assert.Equal(t, username, result.Username)
	assert.False(t, result.IsAdmin)
	assert.False(t, result.Verified)
}
