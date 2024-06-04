package service

import "store.api/dto"

type UserService interface {
	ById(id uint) (*dto.PrivateUserInfo, error)
}
