package controllers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/auth"
	"store.api/config"
	"store.api/controller"
	"store.api/dto"
	"store.api/service"
)

func createAuthController(service service.UserService) *controller.AuthController {
	return controller.NewAuthController(
		service,
		auth.NewJwtMiddleware(&config.Configuration{
			AuthKey: "test secret key",
		}, service).Middle.LoginHandler,
	)
}

func Test_ShouldRegister(t *testing.T) {
	// arrange
	service := createMockUserService()
	controller := createAuthController(service)
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    "mail@mail.com",
	}
	service.On("Register", mock.Anything).Return(nil)

	c, w := createTestContext(data)

	// act
	controller.Register(c)

	// assert
	assert.Equal(t, w.Code, 200)
}

func Test_ShouldNotRegister(t *testing.T) {
	// arrange
	service := createMockUserService()
	controller := createAuthController(service)
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    "mail@mail.com",
	}
	service.On("Register", mock.Anything).Return(errors.New(""))

	c, w := createTestContext(data)

	// act
	controller.Register(c)

	// assert
	assert.Equal(t, w.Code, 400)
}

func Test_ShouldNotRegisterBadData(t *testing.T) {
	// arrange
	service := createMockUserService()
	controller := createAuthController(service)
	data := []string{"first", "second"}
	c, w := createTestContext(data)

	// act
	controller.Register(c)

	// assert
	assert.Equal(t, w.Code, 400)
}

func Test_ShouldLogin(t *testing.T) {
	// arrange
	service := createMockUserService()
	controller := createAuthController(service)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}
	service.On("Login", mock.Anything).Return(&dto.PrivateUserInfo{
		Username: "user",
	}, nil)

	c, w := createTestContext(data)

	// act
	controller.Login(c)

	// assert
	assert.Equal(t, w.Code, 200)
}

func Test_ShouldNotLogin(t *testing.T) {
	// arrange
	service := createMockUserService()
	controller := createAuthController(service)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}
	service.On("Login", mock.Anything).Return(nil, errors.New(""))

	c, w := createTestContext(data)

	// act
	controller.Login(c)

	// assert
	assert.Equal(t, w.Code, 401)
}
