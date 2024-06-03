package repository

import "store.api/model"

type LanguageRepository interface {
	All() []*model.Language
}
