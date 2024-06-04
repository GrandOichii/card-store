package dto

import (
	"strconv"

	"store.api/model"
)

type PrivateUserInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

func NewPrivateUserInfo(user *model.User) *PrivateUserInfo {
	result := PrivateUserInfo{
		Id:       strconv.FormatUint(uint64(user.ID), 10),
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	return &result
}
