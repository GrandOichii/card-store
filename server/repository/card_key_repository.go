package repository

import "store.api/model"

type CardKeyRepository interface {
	All() []*model.CardKey
}
