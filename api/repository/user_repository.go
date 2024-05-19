package repository

import "store.api/model"

type UserRepository interface {
	Save(*model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
}
