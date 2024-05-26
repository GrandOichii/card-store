package dto

import "store.api/model"

type PostCollection struct {
	Name        string `json:"name" validate:"required,gte=3"`
	Description string `json:"description"`
}

func (c *PostCollection) ToCollection() *model.Collection {
	return &model.Collection{
		Name:        c.Name,
		Description: c.Description,
		Cards:       []model.CollectionSlot{},
	}
}
