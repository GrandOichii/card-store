package endpoint_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
)

func Test_ShouldFetchAll(t *testing.T) {
	// arrange
	r, _ := setupRouter()

	// act
	w, _ := req(r, t, "GET", "/api/v1/card", nil, "")

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_ShouldNotCreate(t *testing.T) {
	// arrange
	r, _ := setupRouter()

	// act
	w, _ := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
		Name: "card name",
		Text: "card text",
	}, "")

	// assert
	assert.Equal(t, 403, w.Code)
}

// func Test_ShouldNotCreateNotEnoughPrivileges(t *testing.T) {
// 	// arrange
// 	r, _ := setupRouter()
// 	createUser(r, t, "user", "password", "mail@mail.com")

// 	// act
// 	w, _ := req(r, t, "POST", "/api/v1/card", dto.CreateCard{
// 		Name: "card name",
// 		Text: "card text",
// 	}, "")

// 	// assert
// 	assert.Equal(t, 403, w.Code)
// }
