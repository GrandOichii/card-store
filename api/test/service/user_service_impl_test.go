package service_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/dto"
	"store.api/model"
	"store.api/security"
	"store.api/service"
)

func createUserService(repo *MockUserRepository) service.UserService {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return service.NewUserServiceImpl(
		repo,
		validate,
	)
}

func Test_ShouldRegister(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createUserService(repo)
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    "mail@mail.com",
	}

	repo.On("Save", mock.Anything).Return(nil)
	repo.On("FindByUsername", data.Username).Return(nil)
	repo.On("FindByEmail", data.Email).Return(nil)

	// act
	err := service.Register(&data)

	// assert
	assert.Nil(t, err)
}

func Test_ShouldNotRegisterUsernameTaken(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createUserService(repo)
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    "mail@mail.com",
	}

	repo.On("FindByUsername", data.Username).Return(&model.User{})
	repo.On("FindByEmail", data.Email).Return(nil)

	// act
	err := service.Register(&data)

	// assert
	assert.NotNil(t, err)
}

func Test_ShouldNotRegisterEmailTaken(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createUserService(repo)
	email := "mail@mail.com"
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    email,
	}

	repo.On("FindByUsername", data.Username).Return(nil)
	repo.On("FindByEmail", data.Email).Return(&model.User{
		Email:    "mail@mail.com",
		Verified: true,
	})

	// act
	err := service.Register(&data)

	// assert
	assert.NotNil(t, err)
}

func Test_ShouldNotLogin(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createUserService(repo)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}

	repo.On("FindByUsername", data.Username).Return(nil)

	// act
	login, err := service.Login(&data)

	// assert
	assert.Nil(t, login)
	assert.NotNil(t, err)
}

func Test_ShouldNotLoginIncorrectPassword(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createUserService(repo)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}
	existing := model.User{
		Username:     data.Username,
		PasswordHash: "passwordHash",
	}

	repo.On("FindByUsername", data.Username).Return(&existing)

	// act
	login, err := service.Login(&data)

	// assert
	assert.Nil(t, login)
	assert.NotNil(t, err)
}

func Test_ShouldLogin(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createUserService(repo)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}
	hash, _ := security.HashPassword(data.Password)
	existing := model.User{
		Username:     data.Username,
		PasswordHash: hash,
	}

	repo.On("FindByUsername", data.Username).Return(&existing)

	// act
	login, err := service.Login(&data)

	// assert
	assert.NotNil(t, login)
	assert.Nil(t, err)
}
