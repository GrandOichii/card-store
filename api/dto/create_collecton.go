package dto

import "store.api/model"

type CreateCollection struct {
	Name        string `json:"name" validate:"required,gt=3"`
	Description string `json:"description"`
}

func (c *CreateCollection) ToCollection() *model.Collection {
	return &model.Collection{
		Name:        c.Name,
		Description: c.Description,
		Cards:       []model.CardSlot{},
	}
}
