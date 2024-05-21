package endpoint_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"store.api/dto"
)

func Test_ShouldRegister(t *testing.T) {
	r, _ := setupRouter()
	w, _ := req(r, t, "POST", "/api/v1/auth/register", dto.RegisterDetails{
		Username: "user1",
		Password: "password",
		Email:    "mail@mail.com",
	}, "")

	assert.Equal(t, 201, w.Code)
}

func Test_ShouldNotRegisterBadRequest(t *testing.T) {
	r, _ := setupRouter()

	testCases := []struct {
		desc string
		user dto.RegisterDetails
	}{
		{
			desc: "No username",
			user: dto.RegisterDetails{
				Username: "",
				Password: "password",
				Email:    "mail@mail.com",
			},
		},
		{
			desc: "Short username",
			user: dto.RegisterDetails{
				Username: "u",
				Password: "password",
				Email:    "mail@mail.com",
			},
		},
		{
			desc: "Long username",
			user: dto.RegisterDetails{
				Username: "usernameusernameusernameusernameusernameusernameusernameusernameusernameusernameusernameusernameusernameusernameusernameusername",
				Password: "password",
				Email:    "mail@mail.com",
			},
		},
		{
			desc: "No password",
			user: dto.RegisterDetails{
				Username: "user",
				Password: "",
				Email:    "mail@mail.com",
			},
		},
		{
			desc: "Short password",
			user: dto.RegisterDetails{
				Username: "user",
				Password: "passwor",
				Email:    "mail@mail.com",
			},
		},
		{
			desc: "Long password",
			user: dto.RegisterDetails{
				Username: "user",
				Password: "passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
				Email:    "mail@mail.com",
			},
		},
		{
			desc: "No email",
			user: dto.RegisterDetails{
				Username: "user",
				Password: "passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
				Email:    "",
			},
		},
		{
			desc: "Invalid email 1",
			user: dto.RegisterDetails{
				Username: "user",
				Password: "passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
				Email:    "mail@mail",
			},
		},
		{
			desc: "Invalid email 2",
			user: dto.RegisterDetails{
				Username: "user",
				Password: "passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
				Email:    "mailmail.ru",
			},
		},
		{
			desc: "Invalid email 3",
			user: dto.RegisterDetails{
				Username: "user",
				Password: "passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
				Email:    "@mail.ru",
			},
		},
		{
			desc: "Invalid email 3",
			user: dto.RegisterDetails{
				Username: "user",
				Password: "passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword",
				Email:    "mail@mail.",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			w, _ := req(r, t, "POST", "/api/v1/auth/register", tC.user, "")

			assert.Equal(t, 400, w.Code)
		})
	}
}

func Test_ShouldNotRegisterUsernameExists(t *testing.T) {
	r, _ := setupRouter()
	data := dto.RegisterDetails{
		Username: "user1",
		Password: "password",
		Email:    "mail@mail.com",
	}
	req(r, t, "POST", "/api/v1/auth/register", data, "")
	w, _ := req(r, t, "POST", "/api/v1/auth/register", data, "")

	assert.Equal(t, 400, w.Code)
}

func Test_ShouldLogin(t *testing.T) {
	r, _ := setupRouter()

	register := dto.RegisterDetails{
		Username: "user1",
		Email:    "mail@mail.com",
		Password: "password",
	}
	login := dto.LoginDetails{
		Username: register.Username,
		Password: register.Password,
	}
	req(r, t, "POST", "/api/v1/auth/register", register, "")
	w, _ := req(r, t, "POST", "/api/v1/auth/login", login, "")

	assert.Equal(t, 200, w.Code)
}

func Test_ShouldNotLoginWrongUsername(t *testing.T) {
	r, _ := setupRouter()

	data := dto.LoginDetails{
		Username: "user1",
		Password: "password",
	}
	w, _ := req(r, t, "POST", "/api/v1/auth/login", data, "")

	assert.Equal(t, 401, w.Code)
}

func Test_ShouldNotLoginWrongPassword(t *testing.T) {
	r, _ := setupRouter()

	register := dto.RegisterDetails{
		Username: "user1",
		Password: "password",
		Email:    "mail@mail.com",
	}

	login := dto.LoginDetails{
		Username: "user1",
		Password: "password1",
	}
	req(r, t, "POST", "/api/v1/auth/register", register, "")
	w, _ := req(r, t, "POST", "/api/v1/auth/login", login, "")

	assert.Equal(t, 401, w.Code)
}
