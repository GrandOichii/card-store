package dto

import (
	"store.api/model"
	"store.api/security"
)

type RegisterDetails struct {
	Username string `json:"username" validate:"required,gte=4,lte=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=20"`
}

func (o *RegisterDetails) ToUser() (*model.User, error) {
	passHash, err := security.HashPassword(o.Password)
	if err != nil {
		return nil, err
	}
	result := model.User{
		Username:     o.Username,
		Email:        o.Email,
		PasswordHash: passHash,
		Verified:     false,
	}

	return &result, nil
}
