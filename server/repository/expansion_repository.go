package repository

import "store.api/model"

type ExpansionRepository interface {
	All() []*model.Expansion
}
