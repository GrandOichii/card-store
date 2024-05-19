package repository

import (
	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

type UserDbRepository struct {
	UserRepository

	db     *gorm.DB
	config *config.Configuration
}

func NewUserDbRepository(db *gorm.DB, config *config.Configuration) *UserDbRepository {
	return &UserDbRepository{
		db:     db,
		config: config,
	}
}

func (r *UserDbRepository) Save(user *model.User) error {
	err := r.db.Create(user).Error
	return err
}

func (r *UserDbRepository) FindByUsername(username string) *model.User {
	var result model.User
	err := r.db.Where("username=?", username).Find(&result).Error
	if err != nil {
		panic(err)
	}

	// TODO is there a better way?
	if result.Username != username {
		return nil
	}

	return &result
}

func (r *UserDbRepository) FindByEmail(email string) *model.User {
	var result model.User
	err := r.db.Where("email=?", email).Find(&result).Error
	if err != nil {
		panic(err)
	}

	// TODO is there a better way?
	if result.Email != email {
		return nil
	}

	return &result
}
