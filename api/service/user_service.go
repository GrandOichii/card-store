package service

import "store.api/dto"

type UserService interface {
	Register(*dto.RegisterDetails) (*dto.PrivateUserInfo, error)
}
